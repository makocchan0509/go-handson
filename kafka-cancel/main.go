package main

import (
	"context"
	"fmt"
	"github.com/Shopify/sarama"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var wg sync.WaitGroup

type exampleConsumerGroupHandler struct{}

var client sarama.Client
var group sarama.ConsumerGroup

func (exampleConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (exampleConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (h exampleConsumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	wg.Add(1)
	for msg := range claim.Messages() {
		fmt.Printf("Message end topic:%q partition:%d offset:%d\n", msg.Topic, msg.Partition, msg.Offset)
		time.Sleep(time.Second * 1)
		sess.MarkMessage(msg, "")
	}
	//time.Sleep(time.Second * 3)
	fmt.Println("Message process end")

	wg.Done()
	return nil
}

func main() {

	signals := []os.Signal{
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	}
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, signals...)

	ctx, cancel := context.WithCancel(context.TODO())
	kafkaConsumeStart(ctx)
	<-sigCh
	fmt.Println("received signal... the process will finish consume")
	fmt.Println("cancel context")
	cancel()
	fmt.Println("Close kafka client & group")
	_ = client.Close()
	_ = group.Close()
	time.Sleep(time.Second * 3)
	wg.Wait()
}

func kafkaConsumeStart(ctx context.Context) {

	// Init config, specify appropriate version
	config := sarama.NewConfig()
	config.Version = sarama.V1_0_0_0
	config.Consumer.Return.Errors = true

	// Start with a client
	var err error
	client, err = sarama.NewClient([]string{"localhost:9092", "localhost:9093", "localhost:9094"}, config)
	if err != nil {
		panic(err)
	}

	// Start a new consumer group
	group, err = sarama.NewConsumerGroupFromClient("consumer-group", client)
	if err != nil {
		panic(err)
	}

	// Track errors
	go func() {
		for err := range group.Errors() {
			fmt.Println("ERROR", err)
			panic(err)
		}
	}()
	topics := []string{"test-topic"}

	go func() {
		for {
			handler := exampleConsumerGroupHandler{}
			err := group.Consume(ctx, topics, handler)
			if err != nil {
				fmt.Println("ERROR", err)
				panic(err)
			}
			select {
			case <-ctx.Done():
				fmt.Println("go routine consume done.")
				return
			default:
				fmt.Println("process progress next loop")
			}
		}
	}()
	return
}

/**
func parent_process(ctx context.Context) {
	num := 1
	depth := 0
	for {
		fmt.Printf("parent process start count: %d \n", num)
		go child_process(ctx, num, depth)
		num += 1
		select {
		case <-ctx.Done():
			fmt.Println("received signal will finish process")
			return
		default:
			time.Sleep(time.Second * 5)
		}
	}
}

func child_process(ctx context.Context, num int, depth int) {
	wg.Add(1)
	childCtx, _ := context.WithCancel(ctx)
	if depth > 5 {
		fmt.Printf("process done count: %d \n", num)
		wg.Done()
		return
	}
	fmt.Printf("child process count: %d depth: %d \n", num, depth)
	depth += 1
	time.Sleep(time.Second * 3)
	child_process(childCtx, num, depth)
}
**/
