package main

import (
	"log"

	"github.com/areknoster/attendgo/photo/linuxcapturer"
)

func main() {
	decoder, err := linuxcapturer.NewYUYVDecoder()
	if err != nil {
		log.Fatal("get YUYV decoder: ", err)
	}
	capturer, err := linuxcapturer.Open(linuxcapturer.FormatYUYV, decoder)
	if err != nil {
		log.Fatal("open linux capturer", err)
	}
	capturer.Capture()
	log.Print("success")
}
