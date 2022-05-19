package testsubscribers

import (
	"fmt"
	"github.com/areknoster/attendgo/domain"
)

type Printer struct{}

func (p Printer) Handle(ev domain.Event) {
	fmt.Println(ev)
}
