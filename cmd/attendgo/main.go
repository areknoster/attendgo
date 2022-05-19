package main

import (
	"context"
	"github.com/areknoster/attendgo/domain/notary"
	"github.com/areknoster/attendgo/id"
	"github.com/areknoster/attendgo/id/keyboard"
	"github.com/areknoster/attendgo/mempubsub"
	"github.com/areknoster/attendgo/memstorage"
	"github.com/areknoster/attendgo/photo"
	"github.com/areknoster/attendgo/photo/fakecapturer"
	"github.com/areknoster/attendgo/test/testsubscribers"
	"log"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	ps := mempubsub.NewPubSub(mempubsub.Config{BufSize: 20})
	defer ps.Close()

	ps.Register(testsubscribers.Printer{})
	ps.Register(id.NewHandler(ps))
	kb := keyboard.NewListener(ps)
	storage := memstorage.New()
	capturer := &fakecapturer.Capturer{}
	photographer, err := photo.NewFacePhotographer(capturer, ps)
	if err != nil {
		log.Fatal(err)
	}
	theNotary := notary.New(ps, storage, photographer)
	ps.Register(theNotary)

	log.Fatal(kb.Run(ctx))
}
