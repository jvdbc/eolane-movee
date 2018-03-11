package main

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"testing"
)

// func Test_main(t *testing.T) {
// 	tests := []struct {
// 		name string
// 	}{
// 	// TODO: Add test cases.
// 	}
// 	for range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			main()
// 		})
// 	}
// }

// func Test_oldMain(t *testing.T) {
// 	type args struct {
// 		frames []string
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 	}{
// 	// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			oldMain(tt.args.frames)
// 		})
// 	}
// }

func errIf(t *testing.T, err error) {
	if err != nil {
		t.Error(err)
	}
}

func Test_truc(t *testing.T) {
	data, err := hex.DecodeString("01aa02")
	errIf(t, err)
	sep, err := hex.DecodeString("AA")
	errIf(t, err)
	r := bytes.Split(data, sep)
	t.Logf("Test_Truc print: %s", "hum ...")
	t.Logf("%#v \n", r)

}

func Test_again(t *testing.T) {
	payload := []byte{0x00, 0x01, 0x02, 0x03, 0x00, 0xE2}
	var nb int16
	buf := bytes.NewReader(payload[4:])
	if err := binary.Read(buf, binary.BigEndian, &nb); err != nil {
		errIf(t, err)
	}
	t.Logf("nb: %d", nb)
}
