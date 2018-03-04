package main

import (
	"bytes"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jvdbc/eolane-movee"

	"github.com/reactivex/rxgo/handlers"

	"github.com/reactivex/rxgo/observable"
)

func unTuple(result interface{}, err error) interface{} {
	if err != nil {
		return err
	}
	return result
}

func onInputs(item interface{}) interface{} {
	if str, ok := item.(string); ok {
		return unTuple(hex.DecodeString(str))
	}

	return fmt.Errorf("Unable to cast into string: %#v of type: %T", item, item)
}

func onFrame(item interface{}) interface{} {
	if frm, ok := item.([]byte); ok {
		return unTuple(frame.Payload(frm).Parse())
	}

	return fmt.Errorf("Unable to cast into []byte: %#v of type: %T", item, item)
}

func printFrame(item interface{}) {
	if frm, ok := item.(frame.MoveeFrame); ok {
		fmt.Printf("%s \n", frm)
		return
	}

	log.Printf("Unable to cast into moveeFrame: %#v of type: %T", item, item)
}

// SliceByte type
type sliceByte struct {
	value [][]byte
	index int
}

// Next contract to rx.Iterator
func (s sliceByte) Next() (interface{}, error) {
	if s.index < len(s.value) {
		return s.value[s.index], nil
	}
	return nil, fmt.Errorf("End of SliceByte: %#v", s)
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

	// inputs, err := iterable.New(interface{}(args))

	// if err != nil {
	// 	log.Fatal(err)
	// }

	var frames [][]byte
	buildFrames := handlers.NextFunc(func(item interface{}) {
		if data, ok := item.([]byte); ok {
			b := bytes.Split(data, []byte(sep))
			for _, r := range b {
				frames = append(frames, r)
			}
		}

		log.Printf("Unable to cast into []byte %#v of type %T", item, item)
	})

	wait := observable.Just(args).Map(onInputs).Subscribe(buildFrames)

	<-wait
	// frameIte := sliceByte{value: frames}
	// ite, err := iterable.New(frames)

	// if err != nil {
	// 	log.Fatal(err)
	// }

	wait = observable.Just(frames).Map(onFrame).Subscribe(buildFrames)
	<-wait

	// wait = observable.From(ite).Map()
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
