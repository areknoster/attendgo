package testsubscribers

import (
	"github.com/areknoster/attendgo/domain"
)

type Accumulator struct {
	Received []domain.Event
}

func (a *Accumulator) Handle(ev domain.Event) {
	a.Received = append(a.Received, ev)
}
