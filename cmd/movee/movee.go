package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/reactivex/rxgo/handlers"

	"github.com/reactivex/rxgo/iterable"
	"github.com/reactivex/rxgo/observable"
)

// MoveeFrameType data type
type MoveeFrameType byte

const (
	// Unknown Not a valid frame
	Unknown MoveeFrameType = 0x00
	// Alive frame
	Alive MoveeFrameType = 0x01
	// Temperature frame
	Temperature MoveeFrameType = 0x02
	// Shock frame
	Shock MoveeFrameType = 0x04
	// Tilt frame
	Tilt MoveeFrameType = 0x08
	// Orient frame
	Orient MoveeFrameType = 0x10
	// Motion frame
	Motion MoveeFrameType = 0x20
	// Activity frame
	Activity MoveeFrameType = 0x40
	// Rotation frame
	Rotation MoveeFrameType = 0x80
	// Vibration frame
	Vibration MoveeFrameType = 0x86
	// Information frame
	Information MoveeFrameType = 0xFE
	// Service frame
	Service MoveeFrameType = 0xFF
)

// Payload type alias on slice of bytes
type Payload []byte

// MoveeFrame type
type MoveeFrame struct {
	BatteryLevel float32
	Temperature  int8
	Type         MoveeFrameType
}

func batteryLevel(b byte) float32 {
	return ((3.6-2.8)/255.0)*float32(b) + 2.8
}

func temperature(b byte) int8 {
	return int8(b)
}

func (f MoveeFrame) print() {
	switch f.Type {
	case Unknown:
	case Alive:
	case Temperature:
	case Shock:
	case Tilt:
	case Orient:
	case Motion:
	case Activity:
	case Rotation:
	case Vibration:
	case Information:
	case Service:
	}
}

// Parse read a slice of data to get the movee frame
func (p Payload) Parse() (MoveeFrame, error) {
	if p == nil {
		return MoveeFrame{}, fmt.Errorf("The payload must not be nil")
	}

	if len(p) < 4 {
		return MoveeFrame{}, fmt.Errorf("The payload must have at least 4 bytes: %x", p)
	}

	result := MoveeFrame{BatteryLevel: batteryLevel(p[0]), Temperature: temperature(p[1])}

	return result, nil
}

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
			payload := Payload(data)

			return unravel(payload.Parse())
		}

		return fmt.Errorf("Unable to cast into []byte: %s", item)
	}

	printScreen := handlers.NextFunc(func(item interface{}) {
		f, ok := item.(MoveeFrame)
		if !ok {
			log.Printf("Unable to cast into into MoveeFrame: %s", item)
		}

		f.print()
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
