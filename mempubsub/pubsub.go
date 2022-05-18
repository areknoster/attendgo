package mempubsub

import (
	"github.com/areknoster/attendgo/domain"
)

type Config struct {
	BufSize uint
}

type PubSub[E domain.Event] struct {
	cfg  Config
	subs []chan E
}

func NewPubSub[E domain.Event](cfg Config) *PubSub[E] {
	return &PubSub[E]{
		cfg: cfg,
	}
}

func (ps *PubSub[E]) Publish(ev E) {
	for _, sub := range ps.subs {
		sub <- ev
	}
}

func (ps *PubSub[E]) Register(s domain.Subscriber[E]) {
	queue := make(chan E, ps.cfg.BufSize)
	ps.subs = append(ps.subs, queue)
	go func() {
		for ev := range queue {
			s.Handle(ev)
		}
	}()
}

func (ps *PubSub[E]) Close() {
	for _, queue := range ps.subs {
		close(queue)
	}
}
