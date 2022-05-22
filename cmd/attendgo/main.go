package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/areknoster/attendgo/domain"
	"github.com/areknoster/attendgo/domain/notary"
	"github.com/areknoster/attendgo/id"
	"github.com/areknoster/attendgo/id/keypad"
	"github.com/areknoster/attendgo/mempubsub"
	"github.com/areknoster/attendgo/memstorage"
	"github.com/areknoster/attendgo/photo"
	"github.com/areknoster/attendgo/photo/linuxcapturer"
	"github.com/areknoster/attendgo/signals"
	"github.com/areknoster/attendgo/signals/gpio"
	"github.com/areknoster/attendgo/signals/memio"
	"github.com/areknoster/attendgo/test/testsubscribers"
	"github.com/warthog618/gpiod/device/rpi"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	ps := mempubsub.NewPubSub(mempubsub.Config{BufSize: 20})
	defer ps.Close()

	ps.Register(testsubscribers.Printer{})
	ps.Register(id.NewHandler(ps))

	kp, closer := GPIO(ps)
	defer closer.Close()

	storage := memstorage.New()
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
	theNotary := notary.New(ps, storage, photographer)
	ps.Register(theNotary)

	ps.Register(testsubscribers.PhotoSaver{Dir: "."})
	log.Fatal(kp.Run(ctx))
}

func GPIO(ps domain.Publisher) (*keypad.Keypad, closer) {
	// 1 and 5 work
	cl := closer{}
	greenDiode := gpio.NewGPIOIn(rpi.GPIO23)
	cl.Add(greenDiode)

	cnf := keypad.Config{
		WaitForSignal:     time.Millisecond,
		SweepInteval:      10 * time.Millisecond,
		WaitBetweenInputs: 300 * time.Millisecond,
		Rows: [4]signals.Input{
			gpio.NewGPIOIn(rpi.GPIO9),
			gpio.NewGPIOIn(rpi.GPIO11),
			memio.NewMemIO(),
			memio.NewMemIO(),
		},
		Columns: [4]signals.Output{
			gpio.NewGPIOOut(rpi.GPIO10),
			gpio.NewGPIOOut(rpi.GPIO12),
			memio.NewMemIO(),
			memio.NewMemIO(),
		},
	}
	return keypad.NewKeypad(ps, cnf), cl
}

type closer []interface{ Close() }

func (c closer) Close() {
	for _, closer := range c {
		closer.Close()
	}
}

func (c closer) Add(cl interface{ Close() }) {
	c = append(c, cl)
}
