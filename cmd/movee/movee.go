package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s: frame1frame2frameN\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(-1)
	}

	args := os.Args[1]

	var sep = "aa"

	// Maybe upper
	if !strings.Contains(args, sep) {
		sep = "AA"
	}

	frames := strings.Split(args, sep)

	for i, frame := range frames {
		var data []byte
		var err error

		// onError := handlers.ErrFunc(func(err error) {
		// 	log.Printf("frame %d: %s is not valid: %s ", i, frame, err)
		// })

		// onHex := handlers.NextFunc(func(item interface{}) {
		// 	data := item.([]byte)
		// 	switch {

		// 	}
		// })

		// popHex := func() interface{} {
		// 	if data, err := hex.DecodeString(frame); err != nil {
		// 		return err
		// 	}
		// 	return data
		// }

		// observer := observer.New(onFrame, onError)

		// observable.Start(popHex).Subscribe()

		// var data, err = hex.DecodeString(frame)

		if data, err = hex.DecodeString(frame); err != nil {
			log.Printf("frame %d: %s is not valid !", i, frame)
			continue
		}

		if len(data) < 3 {
			log.Printf("frame %d: %s has length < 3", i, frame)
			continue
		}

		fmt.Printf("frame %d: %s \n", i, data)

	}
}
