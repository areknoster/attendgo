package keypad

import (
	"context"
	"fmt"
	"time"

	"github.com/areknoster/attendgo/domain"
	"github.com/areknoster/attendgo/signals"
)

type (
	Rows    [4]signals.Input
	Columns [4]signals.Output
)

type Config struct {
	WaitForSignal     time.Duration
	SweepInteval      time.Duration
	WaitBetweenInputs time.Duration
	Rows              Rows
	Columns           Columns
}

func New(pub domain.Publisher, config Config) *Keypad {
	return &Keypad{
		pub:    pub,
		config: config,
	}

}

type Keypad struct {
	pub    domain.Publisher
	config Config
}

func (k *Keypad) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("keypad context cancelled: %w", ctx.Err())
		case <-time.After(k.config.SweepInteval):
			glyph, ok := k.sweep()
			if !ok {
				continue
			}
			k.pub.Publish(domain.EventKeyClicked{Glyph: glyph})
			time.Sleep(k.config.WaitBetweenInputs)
		}

	}
}

func (k *Keypad) sweep() (glyph rune, ok bool) {
	for col, out := range k.config.Columns {
		out.Set(true)
		time.Sleep(k.config.WaitForSignal)
		for row, in := range k.config.Rows {
			if in.Value() {
				out.Set(false)
				return keysMapping[input{row, col}], true
			}
		}
		out.Set(false)
	}
	return 0, false
}
