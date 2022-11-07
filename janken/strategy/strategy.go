package strategy

import (
	"math/rand"
	"time"
)

type Strategy interface {
	NextHand() *Hand
	study(win bool)
}
type WinningStrategy struct {
	won      bool
	prevHand *Hand
}

func NewWinningStrategy() *WinningStrategy {
	return &WinningStrategy{
		won:      false,
		prevHand: nil,
	}
}

func (w *WinningStrategy) NextHand() *Hand {
	if !w.won {
		rand.Seed(time.Now().UnixNano())
		w.prevHand = getHand(rand.Intn(3))
	}
	return w.prevHand
}

func (w *WinningStrategy) study(win bool) {
	w.won = win
}

type CircularStrategy struct {
	hand *Hand
}

func NewCircularStrategy() *CircularStrategy {
	return &CircularStrategy{
		hand: &Hand{
			handValue: 0,
		},
	}
}

func (c *CircularStrategy) NextHand() *Hand {
	return getHand(c.hand.handValue)
}

func (c *CircularStrategy) study(win bool) {
	c.hand.handValue = (c.hand.handValue + 1) % 3
}
