package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/validator"
	"github.com/oklog/ulid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/net/context"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

//TODO tracing
const (
	SERVICE_NAME = "trace-demo"
	ENVIRONMENT  = "local"
)
const logKey string = "log"

type User struct {
	UserId      string `json:"user_id,omitempty" validate:"required,min=1,max=4"`
	Name        string `json:"name,omitempty" validate:"required"`
	Password    string `json:"password,omitempty" validate:"required,min=1,max=10"`
	Description string `json:"description,omitempty"`
}

func (u *User) create(conn *dbConn) error {
	if err := conn.Exec().Create(u).Error; err != nil {
		fmt.Errorf("user create error: %w", err)
		return err
	}
	return nil
}

type responseWriterWrapper struct {
	status int
	http.ResponseWriter
	start time.Time
}

func NewResponseWriterWrapper(w http.ResponseWriter) *responseWriterWrapper {
	return &responseWriterWrapper{http.StatusOK, w, time.Now()}
}

func (rw *responseWriterWrapper) Header() http.Header {
	return rw.ResponseWriter.Header()
}

func (rw *responseWriterWrapper) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func getLogger(ctx context.Context) zerolog.Logger {
	return ctx.Value(logKey).(zerolog.Logger)
}

func withLoggerMiddleware(logger zerolog.Logger, conn *dbConn) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			traceId := r.Header.Get("X-Trace-Id")
			if traceId == "" {
				t := time.Now()
				entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
				id := ulid.MustNew(ulid.Timestamp(t), entropy)
				traceId = id.String()
			}
			//			logger = logger.With().Str("traceId", traceId).Logger()
			ctx := context.WithValue(r.Context(), logKey, logger)
			ctx = context.WithValue(ctx, "traceId", traceId)
			w.Header().Set("X-Trace-id", traceId)
			rw := NewResponseWriterWrapper(w)

			conn.Begin()
			defer func() {
				if r := recover(); r != nil {
					conn.RollBack()
				}
			}()
			ctx = context.WithValue(ctx, "conn", conn)

			next.ServeHTTP(rw, r.WithContext(ctx))
			status := rw.status
			if status == 200 {
				if err := conn.Commit(); err != nil {
					conn.RollBack()
				}
			} else {
				conn.RollBack()
			}
			logger.Info().Int("status", status).Str("latency", time.Now().Sub(rw.start).String()).Send()
		})
	}
}

type dbConn struct {
	db *gorm.DB
	tx *gorm.DB
}

func NewDbConn() *dbConn {
	dsn := "user:password@tcp(localhost:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{SkipDefaultTransaction: true})
	if err != nil {
		fmt.Errorf("db connection error: %w", err)
		return nil
	}
	return &dbConn{
		db,
		nil,
	}
}

func getConnFromCtx(c context.Context) *dbConn {
	return c.Value("conn").(*dbConn)
}
func (c *dbConn) Exec() *gorm.DB {
	return c.tx
}

func (c *dbConn) Begin() {
	c.tx = c.db.Begin()
}

func (c *dbConn) Commit() error {
	return c.tx.Commit().Error
}

func (c *dbConn) RollBack() {
	c.tx.Rollback()
}

func main() {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	tp, err := newJaegarProvider("http://localhost:14268/api/traces")
	if err != nil {
		log.Error().Msgf("OpenTelemetry init error %s", err)
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Error().Msgf("Error shutdown OpenTelemetry provider %s", err)
		}
	}()

	conn := NewDbConn()
	//r := mux.NewRouter()
	//r.Use(otelmux.Middleware("my-server"))
	http.Handle("/users/create", withLoggerMiddleware(logger, conn)(http.HandlerFunc(handleCreateUser)))
	//http.Handle("/", r)
	log.Info().Msg("server startup 8080 port")
	http.ListenAndServe(":8080", nil)
}

var tracer = otel.Tracer("mux-server")

func handleCreateUser(w http.ResponseWriter, r *http.Request) {

	traceId := r.Context().Value("traceId").(string)
	ctx, span := tracer.Start(r.Context(), "createUser", trace.WithAttributes(attribute.String("traceId", traceId)))
	defer span.End()

	logger := getLogger(ctx)
	logger = logger.With().Str("traceId", traceId).Logger()
	time.Sleep(time.Millisecond)

	if r.Method != http.MethodPost {
		http.Error(w, `{"detail": "permit only POST method"}`, http.StatusMethodNotAllowed)
		return
	}

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, fmt.Sprintf(`{"detail":"%s"}`, err), http.StatusInternalServerError)
		return
	}
	logger.Info().Str("user_id", user.UserId).Str("name", user.Name).Str("password", user.Password).Str("description", user.Description).Send()

	validate := validator.New()
	if err := validate.Struct(user); err != nil {
		var out []string
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			for _, fe := range ve {
				switch fe.Field() {
				case "UserId":
					out = append(out, fmt.Sprintf("UserId is required and string length min:1 max:4"))
				case "Name":
					out = append(out, fmt.Sprintf("Name is required"))
				case "Password":
					out = append(out, fmt.Sprintf("Password is required and string length min:1 max:10"))
				}
			}
		}
		http.Error(w, fmt.Sprintf(`{"detail":"%s"}`, strings.Join(out, ",")), http.StatusBadRequest)
		return
	}
	conn := getConnFromCtx(ctx)
	if err := user.create(conn); err != nil {
		http.Error(w, fmt.Sprintf(`{"detail":"%s"}`, err), http.StatusInternalServerError)
		return
	}
	dummy(ctx, 0)

	res, err := json.Marshal(user)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"detail":"%s"}`, err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

func newJaegarProvider(url string) (*tracesdk.TracerProvider, error) {
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		return nil, err
	}
	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String(SERVICE_NAME),
				attribute.String("environment", ENVIRONMENT),
			)))
	otel.SetTracerProvider(tp)
	return tp, nil
}

func dummy(ctx context.Context, i int) {
	ctx, childSpan := tracer.Start(ctx, "dummy func;"+strconv.Itoa(i))
	defer childSpan.End()

	time.Sleep(time.Second)

	if i == 10 {
		return
	} else {
		dummy(ctx, i+1)
	}

}
