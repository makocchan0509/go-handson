package hello

import (
	"fmt"
	"testing"
)

type dummyHello struct{}

func (d *dummyHello) doGreeting(msg string) {
	msg += "I'm test code"
	fmt.Println(msg)
}

func TestGreeting(t *testing.T) {
	c := &Client{&dummyHello{}}
	c.People.doGreeting("hello")
}
