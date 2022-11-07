package member

import "time"

type Member interface {
	GetMember() Mem
	getPoint() int
	AddPoint(amount int) int
	SpendPoint(point int) int
	CheckExpire()
	updateTime()
}

type Mem struct {
	Name     string
	Point    int
	Rank     string
	CreateAt time.Time
	UpdateAt time.Time
	ExpireAt time.Time
}

func NewMember(n string, id string) Member {
	var r string
	m := Mem{
		Name:     n,
		Point:    100,
		Rank:     r,
		CreateAt: time.Now(),
		UpdateAt: time.Now(),
		ExpireAt: time.Now().Add(148 * time.Hour),
	}
	if id == "G" {
		m.Rank = "Gold"
		return &Gold{member: m}
	} else {
		m.Rank = "Silver"
		return &Silver{member: m}
	}
}
