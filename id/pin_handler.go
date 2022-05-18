package id

import (
	"github.com/areknoster/attendgo/domain"
	"strings"
)

var _ domain.Subscriber[domain.EventKeyClicked] = (*Handler)(nil)

type Handler struct {
	sb     strings.Builder
	errPub domain.Publisher[domain.EventError]
	idPub  domain.Publisher[domain.EventIDInput]
}

func NewHandler(idPub domain.Publisher[domain.EventIDInput], errPub domain.Publisher[domain.EventError]) *Handler {
	return &Handler{
		errPub: errPub,
		idPub:  idPub,
	}
}

const minPinLength = 4

func (s *Handler) Handle(ev domain.EventKeyClicked) {
	switch ev.Glyph {
	case '*':
		if s.sb.Len() < minPinLength {
			s.errPub.Publish(domain.ErrPinTooShort)
		}

	case '#':
		s.sb.Reset()

	default:
		s.sb.WriteRune(ev.Glyph)
	}
}
