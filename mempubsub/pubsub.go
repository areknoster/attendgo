package mempubsub

import (
	"github.com/areknoster/attendgo/domain"
)

type Config struct {
	BufSize uint
}

type PubSub struct {
	cfg  Config
	subs []chan domain.Event
}

func NewPubSub(cfg Config) *PubSub {
	return &PubSub{
		cfg: cfg,
	}
}

func (ps *PubSub) Publish(ev domain.Event) {
	for _, sub := range ps.subs {
		sub <- ev
	}
}

func (ps *PubSub) Register(s domain.Subscriber) {
	queue := make(chan domain.Event, ps.cfg.BufSize)
	ps.subs = append(ps.subs, queue)
	go func() {
		for ev := range queue {
			s.Handle(ev)
		}
	}()
}

func (ps *PubSub) Close() {
	for _, queue := range ps.subs {
		close(queue)
	}
}
