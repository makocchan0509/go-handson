package hello

import "fmt"

type Hello interface {
	doGreeting(string)
}

type Client struct {
	People Hello
}

type Human struct{}

func (h *Human) doGreeting(msg string) {
	msg += "I'm test target code"
	fmt.Println(msg)
}
func (c *Client) calling() {
	c.People.doGreeting("hello")
}
