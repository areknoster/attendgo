package kaypad

import (
	"context"
	"fmt"
	"github.com/areknoster/attendgo/domain"
	"github.com/areknoster/attendgo/signals"
	"time"
)

type (
	Rows    [4]signals.Input
	Columns [4]signals.Output
)

type Config struct {
	waitForSignal     time.Duration
	sweepInteval      time.Duration
	waitBetweenInputs time.Duration
	rows              Rows
	columns           Columns
}

func NewKeypad(pub domain.Publisher, config Config) *Keypad {
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
		case <-time.After(k.config.sweepInteval):
			glyph, ok := k.sweep()
			if !ok {
				continue
			}
			k.pub.Publish(domain.EventKeyClicked{Glyph: glyph})
			time.Sleep(k.config.waitBetweenInputs)
		}

	}
}

func (k *Keypad) sweep() (glyph rune, ok bool) {
	for col, out := range k.config.columns {
		out.Set(true)
		time.Sleep(k.config.waitForSignal)
		for row, in := range k.config.rows {
			if in.Value() {
				out.Set(false)
				return keysMapping[input{row, col}], true
			}
		}
		out.Set(false)
	}
	return 0, false
}
