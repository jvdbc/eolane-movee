package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/reactivex/rxgo/iterable"
	"github.com/reactivex/rxgo/observable"
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

	frames, err := iterable.New(strings.Split(args, sep))

	if err != nil {
		log.Fatal(err)
	}

	onFrame := func(item interface{}) interface{} {
		if str, ok := item.(string); ok {
			data, err := hex.DecodeString(str)
			if err != nil {
				return err
			}
			return data
		}

		return fmt.Errorf("Unable to get frame %s", item)
	}

	observable.
		From(frames).
		Map(onFrame)

	// connectable.New()

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

}

func oldMain(frames []string) {
	for i, frame := range frames {
		var data []byte
		var err error

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
