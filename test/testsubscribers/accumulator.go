package testsubscribers

import (
	"github.com/areknoster/attendgo/domain"
)

type Accumulator[E domain.Event] struct {
	Received []E
}

func (a *Accumulator[E]) Handle(ev E) {
	a.Received = append(a.Received, ev)
}
