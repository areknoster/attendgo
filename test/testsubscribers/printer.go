package testsubscribers

import (
	"log"

	"github.com/areknoster/attendgo/domain"
)

type Printer struct{}

func (p Printer) Handle(ev domain.Event) {
	log.Println(ev)
}
