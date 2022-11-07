package member

import (
	"fmt"
	"time"
)

type Gold struct {
	member Mem
}

func (g *Gold) GetMember() Mem {
	return g.member
}

func (g *Gold) getPoint() int {
	return g.member.Point
}

func (g *Gold) AddPoint(amount int) int {
	rate := 0.04
	g.member.Point += int(float64(amount) * rate)
	return g.getPoint()
}

func (g *Gold) SpendPoint(point int) int {
	g.member.Point -= point
	g.updateTime()
	return g.getPoint()
}

func (g *Gold) CheckExpire() {
	if g.member.ExpireAt.Before(time.Now()) {
		fmt.Println("Expire!!")
		g.member.Point = 0
	} else {
		fmt.Println("not expire")
	}
}

func (g *Gold) updateTime() {
	t := time.Now()
	g.member.UpdateAt = t
	g.member.ExpireAt = t.Add(148 * time.Hour)
}
