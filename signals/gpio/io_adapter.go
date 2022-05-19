package gpio

import (
	"github.com/warthog618/gpiod"
	"log"
)

type GPIOIn struct {
	line *gpiod.Line
}

func (g *GPIOIn) Close() {
	err := g.line.Close()
	if err != nil {
		log.Print("close GPIO", err.Error())
	}
}

func (g *GPIOIn) Value() bool {
	v, err := g.line.Value()
	if err != nil {
		log.Print("couldn't read input", err.Error())
	}
	return v > 0
}

func NewGPIOIn(gpioNo int) *GPIOIn {
	in, err := gpiod.RequestLine("gpiochip0", gpioNo, gpiod.AsInput)
	if err != nil {
		panic(err)
	}
	return &GPIOIn{
		line: in,
	}
}

type GPIOOut struct {
	line *gpiod.Line
}

func (g *GPIOOut) Close() {
	err := g.line.Close()
	if err != nil {
		log.Print("close GPIO", err.Error())
	}
}

func (g *GPIOOut) Set(v bool) {
	var setVal int
	if v {
		setVal = 1
	}
	err := g.line.SetValue(setVal)
	if err != nil {
		log.Print("couldn't read input", err.Error())
	}
}

func NewGPIOOut(gpioNo int) *GPIOOut {
	out, err := gpiod.RequestLine("gpiochip0", gpioNo, gpiod.AsOutput())
	if err != nil {
		panic(err)
	}
	return &GPIOOut{
		line: out,
	}
}
