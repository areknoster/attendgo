package notary

import (
	"github.com/areknoster/attendgo/domain"
)

type Notary struct {
	errPub domain.Publisher[domain.EventError]
}

var (
	_ domain.Subscriber[domain.EventIDInput] = (*Notary)(nil)
)

func (n *Notary) Handle(ev domain.EventIDInput) {
	//TODO implement me
	panic("implement me")
}

func (n *Notary) Handle(ev domain.EventFacePhotoTaken) {
	//TODO implement me
	panic("implement me")
}
