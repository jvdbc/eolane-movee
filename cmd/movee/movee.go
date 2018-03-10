package main

import (
	"bytes"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"os"

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

func hexString(item interface{}) interface{} {
	if str, ok := item.(string); ok {
		return unTuple(hex.DecodeString(str))
	}
	return fmt.Errorf("Unable to cast into string: %#v of type: %T", item, item)
}

func parse(item interface{}) interface{} {
	if payload, ok := item.([]byte); ok {
		return unTuple(frame.Payload(payload).Parse())
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

// byteIterator type
type byteIterator struct {
	value [][]byte
	index int
}

// Next contract to rx.Iterator
func (s *byteIterator) Next() (interface{}, error) {
	if s.index < len(s.value) {
		res := s.value[s.index]
		s.index++
		return res, nil
	}
	return nil, fmt.Errorf("End of IteByte: %#v", s)
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

	input := os.Args[1]

	var sep = []byte{0xaa}

	var frames [][]byte
	splitFrames := handlers.NextFunc(func(item interface{}) {
		data, ok := item.([]byte)
		if !ok {
			log.Printf("Unable to cast into []byte %#v of type %T", item, item)
			return
		}
		b := bytes.Split(data, sep)
		for _, r := range b {
			frames = append(frames, r)
		}
	})

	wait := observable.
		Just(input).
		Map(hexString).
		Subscribe(splitFrames)

	<-wait

	wait = observable.
		From(&byteIterator{value: frames}).
		Map(parse).
		Subscribe(handlers.NextFunc(printFrame))

	<-wait
}
