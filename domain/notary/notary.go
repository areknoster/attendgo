package notary

import (
	"errors"
	"github.com/areknoster/attendgo/domain"
	"log"
	"sync/atomic"
)

type state uint32

const (
	stateListeningForPins state = iota + 1
	stateTakingPhotos
)

func New(pub domain.Publisher, storage domain.AtendeeStorage, photoCapturer domain.Capturer) *Notary {
	n := &Notary{
		pub:           pub,
		storage:       storage,
		photoCapturer: photoCapturer,
	}
	n.setState(stateListeningForPins)
	return n
}

type Notary struct {
	pub           domain.Publisher
	state         uint32
	attendee      domain.Attendee
	storage       domain.AtendeeStorage
	photoCapturer domain.Capturer
}

var (
	_ domain.Subscriber = (*Notary)(nil)
)

func (n *Notary) Handle(ev domain.Event) {
	switch ev := ev.(type) {
	case domain.EventError:
		switch ev {
		case domain.ErrFacePhotoNotTaken:
			n.attendee = domain.Attendee{}
			n.setState(stateListeningForPins)
		}
	case domain.EventFacePhotoTaken:
		if n.getState() != stateTakingPhotos {
			log.Fatal("EventFacePhotoTaken received in wrong state ")
		}
		n.attendee.Photo = ev.Photo
		err := n.storage.Create(n.attendee)
		if errors.Is(err, domain.ErrStorageAttendeeAlreadyExits) {
			n.pub.Publish(domain.ErrAttendeeAlreadyExists)
			return
		}
		if err != nil {
			log.Fatal("unknown storage error")
		}
		n.setState(stateListeningForPins)
		n.pub.Publish(domain.AtendeeRegisteredEvent{
			Atendee: n.attendee,
		})
	case domain.EventIDInput:
		if n.getState() != stateListeningForPins {
			n.pub.Publish(domain.ErrInputDuringPhotoSession)
		}
		n.attendee.ID = ev.ID
		n.setState(stateTakingPhotos)
		go n.photoCapturer.StartCapturing()
	}
}

func (n *Notary) setState(state state) {
	atomic.StoreUint32(&n.state, uint32(state))
}

func (n *Notary) getState() state {
	return state(atomic.LoadUint32(&n.state))
}
