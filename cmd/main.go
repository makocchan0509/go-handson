package main

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"
)

var wg sync.WaitGroup

func main() {

	b := []byte(strconv.Itoa(5))
	fmt.Println(b)

	fmt.Println(string(b))
	fmt.Println(strconv.Atoi(string(b)))

	//ctx, cancel := context.WithCancel(context.Background())
	//queue := make(chan string)
	//for i := 0; i < 3; i++ {
	//	wg.Add(1)
	//	go fetch(ctx, queue, i)
	//}
	//for i := 0; i < 30; i++ {
	//	queue <- "Hi " + strconv.Itoa(i) + " Mase"
	//}
	//cancel()
	//wg.Wait()
}

func fetch(ctx context.Context, queue chan string, n int) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Exit")
			wg.Done()
			return
		case url := <-queue:
			fmt.Println(url, n)
		}
	}
}

func sleep() {
	defer fmt.Println("done sleep")
	for {
		time.Sleep(1 * time.Second)
	}
}

func count(x int, n int, c chan int) {
	for i := 0; i < x; i++ {
		c <- i
		time.Sleep(time.Duration(n) * time.Millisecond)
	}
}

func channel() {
	c := make(chan int)
	go total(c)
	c <- 10
	fmt.Println("total:", <-c)
}

func total(c chan int) {
	n := <-c
	fmt.Println("n --> ", n)
	t := 0
	for i := 0; i < n; i++ {
		t += i
	}
	c <- t
}

func thead() {
	msg := "start"
	primsg := func(s string, n int) {
		fmt.Println(s, msg)
		time.Sleep(time.Duration(n) * time.Millisecond)
	}

	hello := func(n int) {
		const nm string = "hello"
		for i := 0; i < 10; i++ {
			msg += " h" + strconv.Itoa(i)
			primsg(nm, 50)
		}
	}

	main := func(n int) {
		const nm string = "main"
		for i := 0; i < 10; i++ {
			msg += " m" + strconv.Itoa(i)
			primsg(nm, 100)
		}
	}

	go hello(50)
	main(100)
}

func say(s string, t int) {
	for i := 0; i < 10; i++ {
		fmt.Printf("<%d %s>", i, s)
		time.Sleep(time.Duration(t) * time.Millisecond)
	}
}

type general interface{}

type Data interface {
	set(nm string, g general) Data
	print()
}

type NData struct {
	Name string
	Data int
}

func (d *NData) set(nm string, g general) Data {
	d.Name = nm
	d.Data = g.(int)
	return d
}

func (d *NData) print() {
	fmt.Printf("<%s> value: %d\n", d.Name, d.Data)
}

type SData struct {
	Name string
	Data string
}

func (s *SData) set(nm string, g general) Data {
	s.Name = nm
	s.Data = g.(string)
	return s
}

func (s *SData) print() {
	fmt.Printf("<%s> value: %s\n", s.Name, s.Data)
}

func push(arr []string, s ...string) []string {
	return append(arr, s...)
}

func pop(arr []string) ([]string, string) {
	return arr[:len(arr)-1], arr[len(arr)-1]
}
