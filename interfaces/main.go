package main

import (
	"fmt"
	"go-handson/interfaces/member"
)

func main() {

	m := member.NewMember("masem", "G")
	fmt.Println(m.GetMember())
	fmt.Println(m.AddPoint(1000))
	fmt.Println(m.SpendPoint(10))
	m.CheckExpire()
	fmt.Println(m.GetMember())

	s := member.NewMember("silver", "S")
	fmt.Println(s.GetMember())
	fmt.Println(s.AddPoint(1000))
	fmt.Println(s.SpendPoint(10))
	s.CheckExpire()
	fmt.Println(s.GetMember())
}
