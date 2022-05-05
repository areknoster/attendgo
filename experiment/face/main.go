package main

import (
	"log"

	"github.com/areknoster/go-face"
)

func main(){
	rec, err := face.NewRecognizer("../../models")
	if err != nil {
		log.Fatal("could not load recognizer: ", err)
	}
	_ = rec

}