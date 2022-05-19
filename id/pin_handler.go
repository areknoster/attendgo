package id

import (
	"github.com/areknoster/attendgo/domain"
	"strings"
)

var _ domain.Subscriber = (*Handler)(nil)

type Handler struct {
	sb  strings.Builder
	pub domain.Publisher
}

func NewHandler(pub domain.Publisher) *Handler {
	return &Handler{
		pub: pub,
	}
}

const minPinLength = 4

func (s *Handler) Handle(ev domain.Event) {
	switch ev := ev.(type) {
	case domain.EventKeyClicked:
		switch ev.Glyph {
		case '*':
			if s.sb.Len() < minPinLength {
				s.pub.Publish(domain.ErrPinTooShort)
				s.sb.Reset()
				return
			}
			s.pub.Publish(domain.EventIDInput{ID: domain.ID(s.sb.String())})
			s.sb.Reset()

		case '#':
			s.sb.Reset()

		default:
			s.sb.WriteRune(ev.Glyph)
		}
	}
}
