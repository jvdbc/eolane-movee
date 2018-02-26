package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jvdbc/eolane-movee"

	"github.com/reactivex/rxgo/handlers"

	"github.com/reactivex/rxgo/iterable"
	"github.com/reactivex/rxgo/observable"
)

func unravel(result interface{}, err error) interface{} {
	if err != nil {
		return err
	}
	return result
}

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

	inputs, err := iterable.New(strings.Split(args, sep))

	if err != nil {
		log.Fatal(err)
	}

	onInput := func(item interface{}) interface{} {
		if str, ok := item.(string); ok {
			return unravel(hex.DecodeString(str))
		}

		return fmt.Errorf("Unable to cast into string %s", item)
	}

	onBytes := func(item interface{}) interface{} {
		if data, ok := item.([]byte); ok {
			payload := frame.Payload(data)

			return unravel(payload.Parse())
		}

		return fmt.Errorf("Unable to cast into []byte: %s", item)
	}

	printScreen := handlers.NextFunc(func(item interface{}) {
		_, ok := item.(frame.Header)
		if !ok {
			log.Printf("Unable to cast into into MoveeFrame: %s", item)
		}

		// frame.Print()
	})

	observable.
		From(inputs).
		Map(onInput).
		Map(onBytes).
		Subscribe(printScreen)

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
