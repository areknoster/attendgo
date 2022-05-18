package main

import (
	"context"
	"github.com/areknoster/attendgo/domain"
	"github.com/areknoster/attendgo/id"
	"github.com/areknoster/attendgo/id/keyboard"
	"github.com/areknoster/attendgo/mempubsub"
	"github.com/areknoster/attendgo/test/testsubscribers"
	"log"
	"os/signal"
	"syscall"
)

func initPubSub[E domain.Event]() domain.PubSub[E] {
	return mempubsub.NewPubSub[E](mempubsub.Config{BufSize: 10})
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	errPs := initPubSub[domain.EventError]()
	errPs.Register(testsubscribers.Printer[domain.EventError]{})

	keysPs := initPubSub[domain.EventKeyClicked]()
	idPs := initPubSub[domain.EventIDInput]()

	keysPs.Register(id.NewHandler(idPs, errPs))
	keysPs.Register(testsubscribers.Printer[domain.EventKeyClicked]{})

	kb := keyboard.NewListener(keysPs)
	log.Fatal(kb.Run(ctx))
}
