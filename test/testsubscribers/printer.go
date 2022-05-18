package testsubscribers

import (
	"fmt"
)

type Printer[T fmt.Stringer] struct{}

func (p Printer[T]) Handle(ev T) {
	fmt.Println(ev.String())
}
