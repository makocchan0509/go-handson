package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/Shopify/sarama"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"
)

// Producer部分のみ抜粋

var (
	// kafkaのアドレス
	bootstrapServers = flag.String("bootstrapServers", "localhost:9092", "kafka address")
)

// SendMessage 送信メッセージ
type SendMessage struct {
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}

func main() {
	flag.Parse()

	if *bootstrapServers == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	brokers := strings.Split(*bootstrapServers, ",")
	config := sarama.NewConfig()

	config.Producer.Return.Errors = true
	config.Producer.Return.Successes = true
	config.Producer.Retry.Max = 3

	producer, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		panic(err)
	}

	defer producer.AsyncClose()

	// プロデューサールーチン
	go func() {
	PRODUCER_FOR:
		for {
			time.Sleep(1000 * time.Millisecond)

			timestamp := time.Now().UnixNano()

			send := &SendMessage{
				Message:   "Hello",
				Timestamp: timestamp,
			}

			jsBytes, err := json.Marshal(send)
			if err != nil {
				panic(err)
			}

			msg := &sarama.ProducerMessage{
				Topic: "test-topic",
				Key:   sarama.StringEncoder(strconv.FormatInt(timestamp, 10)),
				Value: sarama.StringEncoder(string(jsBytes)),
			}

			producer.Input() <- msg

			select {
			case <-producer.Successes():
				fmt.Println(fmt.Sprintf("success send. message: %s, timestamp: %d", send.Message, send.Timestamp))
			case err := <-producer.Errors():
				fmt.Println(fmt.Sprintf("fail send. reason: %v", err.Msg))
			case <-ctx.Done():
				break PRODUCER_FOR
			}
		}
	}()

	fmt.Println("go-kafka-example start.")

	<-signals

	fmt.Println("go-kafka-example stop.")
}
