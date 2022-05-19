package signals

import (
	"github.com/areknoster/attendgo/domain"
	"time"
)

type Input interface {
	Value() bool
}

type Output interface {
	Set(bool)
}

type SignalOutputs struct {
	Good    Output
	Warning Output
	Bad     Output
}

var _ domain.Subscriber = (*SignalHandler)(nil)

type SignalHandler struct {
	outputs SignalOutputs
}

func (s *SignalHandler) Handle(ev domain.Event) {
	switch ev := ev.(type) {
	case domain.EventError:
		switch ev {
		case domain.ErrFacePhotoNotTaken:
			s.bad(3)
		case domain.ErrNoFace:
			s.warning(1)
		case domain.ErrInputDuringPhotoSession:
			s.bad(1)
		case domain.ErrTooManyFaces:
			s.warning(2)
		case domain.ErrPinTooShort:
			s.bad(3)
		case domain.ErrAttendeeAlreadyExists:
			s.warning(3)
		}
	case domain.EventIDInput:
		s.good(3)
	case domain.AtendeeRegisteredEvent:
		s.good(3)
	}
}

func (s *SignalHandler) wait() {
	time.Sleep(time.Second / 2)
}

func (s *SignalHandler) repeatOut(repeat int, out Output) {
	go func() {
		for i := 0; i < repeat; i++ {
			out.Set(true)
			s.wait()
			out.Set(false)
			s.wait()
		}
	}()
}

func (s *SignalHandler) good(repeat int) {
	s.repeatOut(repeat, s.outputs.Good)
}

func (s *SignalHandler) warning(repeat int) {
	s.repeatOut(repeat, s.outputs.Warning)
}

func (s *SignalHandler) bad(repeat int) {
	s.repeatOut(repeat, s.outputs.Warning)
}
