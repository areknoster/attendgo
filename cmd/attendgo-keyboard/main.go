package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/areknoster/attendgo/domain"
	"github.com/areknoster/attendgo/domain/notary"
	"github.com/areknoster/attendgo/id"
	"github.com/areknoster/attendgo/id/keyboard"
	"github.com/areknoster/attendgo/mempubsub"
	"github.com/areknoster/attendgo/memstorage"
	"github.com/areknoster/attendgo/photo"
	"github.com/areknoster/attendgo/photo/linuxcapturer"
	photosaver "github.com/areknoster/attendgo/photo/saver"
	"github.com/areknoster/attendgo/signals"
	"github.com/areknoster/attendgo/test/testsubscribers"
	"github.com/areknoster/attendgo/webui"
	"github.com/sethvargo/go-envconfig"
	"golang.org/x/sync/errgroup"
)

type Config struct {
	webui.ServerConfig
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	cfg := Config{}
	handleFatal(envconfig.Process(ctx, &cfg))

	ps := mempubsub.NewPubSub(mempubsub.Config{BufSize: 20})
	defer ps.Close()

	ps.Register(testsubscribers.Printer{})
	ps.Register(id.NewHandler(ps))

	kb := keyboard.New(ps)

	photoStorage := memstorage.NewPhotoStorage()
	ps.Register(photosaver.New(photoStorage))

	decoder, err := linuxcapturer.NewYUYVDecoder()
	if err != nil {
		log.Fatal("get YUYV decoder: ", err)
	}
	capturer, err := linuxcapturer.Open(linuxcapturer.FormatYUYV, decoder)
	if err != nil {
		log.Fatal("open linux capturer", err)
	}
	photographer, err := photo.NewFacePhotographer(capturer, ps)
	if err != nil {
		log.Fatal(err)
	}

	attendeeStorage := memstorage.NewAtendeeStorage()
	nt := notary.New(ps, attendeeStorage, photographer)
	ps.Register(nt)

	ps.Register(initSignalHandler())

	uiHandler, err := webui.NewUIHandler(photoStorage)
	handleFatal(err)
	srv := webui.NewServer(cfg.ServerConfig, uiHandler, webui.NewAttendeesHandler(attendeeStorage))

	eg, ctx := errgroup.WithContext(ctx)
	withCtx := func(r domain.Runner) func() error {
		return func() error {
			return r.Run(ctx)
		}
	}
	eg.Go(withCtx(kb))
	eg.Go(withCtx(srv))

	handleFatal(eg.Wait())
}

func handleFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func initSignalHandler() *signals.SignalHandler {
	return signals.NewSignalHandler(signals.SignalOutputs{
		Good:    signals.NamedPrinterDecorator{Name: "Good"},
		Warning: signals.NamedPrinterDecorator{Name: "Warning"},
		Bad:     signals.NamedPrinterDecorator{Name: "Bad"},
	})
}
