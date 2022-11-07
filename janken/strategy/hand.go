package strategy

const (
	GUU = iota
	CHOO
	PAA
)

type Hand struct {
	handValue int
}

var hands []*Hand

func init() {
	hands = []*Hand{
		&Hand{GUU},
		&Hand{CHOO},
		&Hand{PAA},
	}
}

func getHand(handValue int) *Hand {
	return hands[handValue]
}

func (h *Hand) IsStrongerThan(opponentHand *Hand) bool {
	return 1 == h.fight(opponentHand)
}

func (h *Hand) IsWeakerThan(opponentHand *Hand) bool {
	return -1 == h.fight(opponentHand)
}

func (h *Hand) fight(opponentHand *Hand) int {
	if h == opponentHand {
		return 0
	} else if (h.handValue+1)%3 == opponentHand.handValue {
		return 1
	} else {
		return -1
	}
}
