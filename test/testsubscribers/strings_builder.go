package testsubscribers

import (
	"github.com/areknoster/attendgo/domain"
	"strings"
)

type StringsBuilder struct {
	sb strings.Builder
}

func (l *StringsBuilder) Handle(ev domain.EventKeyClicked) {
	l.sb.WriteRune(ev.Glyph)
}
func (l *StringsBuilder) String() string {
	return l.sb.String()
}
