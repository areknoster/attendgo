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
	photosaver "github.com/areknoster/attendgo/photo/saver"
	"github.com/areknoster/attendgo/signals"
	"github.com/areknoster/attendgo/signals/gpio"
	"github.com/areknoster/attendgo/test/testsubscribers"
	"github.com/areknoster/attendgo/webui"
	"github.com/sethvargo/go-envconfig"
	"github.com/warthog618/gpiod/device/rpi"
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

	kb, closeKeyboard := InitKeyboard(ps)
	defer closeKeyboard.Close()

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

	ps.Register(InitSignalHandler())

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

func InitSignalHandler() *signals.SignalHandler {
	return signals.NewSignalHandler(signals.SignalOutputs{
		Good:    gpio.NewGPIOOut(rpi.GPIO22),
		Warning: signals.NamedPrinterDecorator{Name: "warning"},
		Bad:     gpio.NewGPIOOut(rpi.GPIO27),
	})
}

func InitKeyboard(ps domain.Publisher) (*keypad.Keypad, closer) {
	// 1 and 5 work
	cl := closer{}
	greenDiode := gpio.NewGPIOIn(rpi.GPIO23)
	cl.Add(greenDiode)

	cnf := keypad.Config{
		WaitForSignal:     time.Millisecond,
		SweepInteval:      10 * time.Millisecond,
		WaitBetweenInputs: 300 * time.Millisecond,
		Rows: [4]signals.Input{
			gpio.NewGPIOIn(rpi.GPIO21),
			gpio.NewGPIOIn(rpi.GPIO9),
			gpio.NewGPIOIn(rpi.GPIO11),
			gpio.NewGPIOIn(rpi.GPIO20),
			// memio.NewMemIO(),
		},
		Columns: [4]signals.Output{
			// signals.NamedPrinterDecorator{},
			// signals.NamedPrinterDecorator{},
			// signals.NamedPrinterDecorator{},
			// signals.NamedPrinterDecorator{},
			gpio.NewGPIOOut(rpi.GPIO26),
			gpio.NewGPIOOut(rpi.GPIO12),
			gpio.NewGPIOOut(rpi.GPIO25),
			gpio.NewGPIOOut(rpi.GPIO8),
		},
	}
	return keypad.New(ps, cnf), cl
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
