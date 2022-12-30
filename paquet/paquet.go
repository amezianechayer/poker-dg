package paquet

import (
	"fmt"
	"math/rand"
	"strconv"
)

type Suit int

func (s Suit) String() string {
	switch s {
	case Spades:
		return "SPADES"
	case Harts:
		return "HARTS"
	case Diamonds:
		return "DIAMONDS"
	case Clubs:
		return "CLUBS"
	default:
		panic("suite de cartes invalides")
	}
}

const (
	Spades Suit = iota
	Harts
	Diamonds
	Clubs
)

type Cartes struct {
	suit  Suit
	value int
}

func (c Cartes) String() string {
	value := strconv.Itoa(c.value)
	if c.value == 1 {
		value = "ACE"
	}
	return fmt.Sprintf("%s of %s %s", value, c.suit, suitToUnicode(c.suit))
}

func NewCartes(s Suit, v int) Cartes {
	if v > 13 {
		panic("La valeur des cartes ne peut pas étre superieure à 13")
	}
	return Cartes{
		suit:  s,
		value: v,
	}

}

type Paquet [52]Cartes

func New() Paquet {
	var (
		nSuits   = 4
		nCartees = 13
		d        = [52]Cartes{}
	)
	x := 0
	for i := 0; i < nSuits; i++ {
		for j := 0; j < nCartees; j++ {
			d[x] = NewCartes(Suit(i), j+1)
			x++

		}

	}

	return shuffle(d)
}

func shuffle(d Paquet) Paquet {
	for i := 0; i < len(d); i++ {
		r := rand.Intn(i + 1)

		if r != i {
			d[i], d[r] = d[r], d[i]
		}

	}

	return d
}

func suitToUnicode(s Suit) string {
	switch s {
	case Spades:
		return "♠"
	case Harts:
		return "♥"
	case Diamonds:
		return "♦"
	case Clubs:
		return "♣"
	default:
		panic("suite de cartes invalides")
	}
}
