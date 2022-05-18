package keyboard

import (
	"context"
	"fmt"
	"github.com/areknoster/attendgo/domain"
	"github.com/eiannone/keyboard"
)

type Listener struct {
	pub domain.Publisher[domain.EventKeyClicked]
}

func NewListener(pub domain.Publisher[domain.EventKeyClicked]) *Listener {
	return &Listener{pub: pub}
}

func (l *Listener) Run(ctx context.Context) error {
	err := keyboard.Open()

	if err != nil {
		return fmt.Errorf("open keyboard: %w", err)
	}
	defer keyboard.Close()
	keys, err := keyboard.GetKeys(10)
	if err != nil {
		return fmt.Errorf("get keys chan: %w", err)
	}
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("keyboard context cancelled: %w", ctx.Err())
		case key, isOpen := <-keys:
			if !isOpen {
				return fmt.Errorf("keyboard closed")
			}
			if key.Key == keyboard.KeyCtrlC {
				return fmt.Errorf("interrupt")
			}
			l.pub.Publish(domain.EventKeyClicked{Glyph: key.Rune})
		}
	}
	return nil
}
