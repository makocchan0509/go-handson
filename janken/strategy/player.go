package strategy

import "fmt"

type Player struct {
	name                           string
	strategy                       Strategy
	wincount, losecount, gamecount int
}

func NewPlayer(nm string, st Strategy) *Player {
	return &Player{
		name:      nm,
		strategy:  st,
		wincount:  0,
		losecount: 0,
		gamecount: 0,
	}
}

func (p *Player) NextHand() *Hand {
	return p.strategy.NextHand()
}

func (p *Player) Win() {
	p.strategy.study(true)
	p.wincount++
	p.gamecount++
}

func (p *Player) Lose() {
	p.strategy.study(false)
	p.losecount++
	p.gamecount++
}

func (p *Player) Even() {
	p.gamecount++
}

func (p *Player) ToString() string {
	str := fmt.Sprintf("[%s: %dgames, %dwin, %dlose]", p.name, p.gamecount, p.wincount, p.losecount)
	return str
}
