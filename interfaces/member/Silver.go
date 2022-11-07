package member

import (
	"fmt"
	"time"
)

type Silver struct {
	member Mem
}

func (g *Silver) GetMember() Mem {
	return g.member
}

func (g *Silver) getPoint() int {
	return g.member.Point
}

func (g *Silver) AddPoint(amount int) int {
	rate := 0.02
	g.member.Point += int(float64(amount) * rate)
	return g.getPoint()
}

func (g *Silver) SpendPoint(point int) int {
	g.member.Point -= point
	g.updateTime()
	return g.getPoint()
}

func (g *Silver) CheckExpire() {
	if g.member.ExpireAt.Before(time.Now()) {
		fmt.Println("Expire!!")
		g.member.Point = 0
	} else {
		fmt.Println("not expire")
	}
}

func (g *Silver) updateTime() {
	t := time.Now()
	g.member.UpdateAt = t
	g.member.ExpireAt = t.Add(148 * time.Hour)
}
