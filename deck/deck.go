package deck

import (
	"fmt"
	"math/rand"
)

type Suit int

func (s Suit) String() string {
	switch s {
	case Spades:
		return "SPADES"
	case Hearts:
		return "HEARTS"
	case Diamonds:
		return "DIAMOND"
	case Clubs:
		return "CLUBS"
	default:
		panic("invalid card Suit")
	}
}

func (c Cards) String() string {
	return fmt.Sprintf("%d of %s %s", c.value, c.Suit, SuitToUnicode(c.Suit))
}

func SuitToUnicode(s Suit) string {
	switch s {
	case Spades:
		return "♠"
	case Hearts:
		return "♥"
	case Diamonds:
		return "♦"
	case Clubs:
		return "♣"
	default:
		panic("invalid card Suit")
	}
}

const (
	Spades Suit = iota
	Hearts
	Diamonds
	Clubs
)

type Cards struct {
	Suit  Suit
	value int
}

func NewCard(s Suit, v int) Cards {
	if v > 13 {
		panic("value cannot be greater than 13")
	}
	return Cards{
		Suit:  s,
		value: v,
	}
}

type Deck [52]Cards

func New() Deck {
	var (
		noSuits = 4
		noCards = 13
		d       = [52]Cards{}
		x       = 0
	)

	for i := 0; i < noSuits; i++ {
		for j := 0; j < noCards; j++ {
			d[x] = NewCard(Suit(i), j+1)
			x++
		}
	}
	return Shuffle(d)
}

func Shuffle(d Deck) Deck {
	for i := 0; i < len(d); i++ {
		r := rand.Intn(i + 1)
		fmt.Print(r)
		if r != i {
			d[i], d[r] = d[r], d[i]
		}
	}
	return d
}
