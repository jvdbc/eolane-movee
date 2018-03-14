package frame

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// MoveeType data type
type MoveeType byte

const (
	// Alive frame
	Alive MoveeType = 0x01
	// Temperature frame
	Temperature MoveeType = 0x02
	// Shock frame
	Shock MoveeType = 0x04
	// Tilt frame
	Tilt MoveeType = 0x08
	// Orient frame
	Orient MoveeType = 0x10
	// Motion frame
	Motion MoveeType = 0x20
	// Activity frame
	Activity MoveeType = 0x40
	// Rotation frame
	Rotation MoveeType = 0x80
	// Vibration frame
	Vibration MoveeType = 0x86
	// Information frame
	Information MoveeType = 0xFE
	// Service frame
	Service MoveeType = 0xFF
)

//go:generate stringer -type=MoveeType

// Payload type alias on slice of bytes
type Payload []byte

// Header type
type Header struct {
	BatteryLevel float32
	Temperature  int8
	FrameType    MoveeType
}

// MoveeFrame interface
type MoveeFrame interface {
}

// AliveFrame type
type AliveFrame struct {
	*Header
}

// TemperatureFrame type
type TemperatureFrame struct {
	*Header
	Temperatures []int8
}

// ShockFrame type
type ShockFrame struct {
	*Header
}

// TiltFrame type
// c1150800000000aac115080000ffecaac115080064ffecaa
type TiltFrame struct {
	*Header
	PitchAngle int16
	RollAngle  int16
}

// OrientFrame type
type OrientFrame struct {
	*Header
	PitchAngle int16
	RollAngle  int16
	YawAngle   int16
}

// MotionFrame type
type MotionFrame struct {
	*Header
	OnMove bool
}

// ActivityFrame type
type ActivityFrame struct {
	*Header
	OnMove   bool
	Duration uint32
}

// RotationFrame type
// c11580000aaa
type RotationFrame struct {
	*Header
	TurnsNumber int16
}

// VibeFrame type
type VibeFrame struct {
	*Header
}

// InformationFrame type
type InformationFrame struct {
	*Header
}

// ServiceFrame type
type ServiceFrame struct {
	*Header
}

func (m Header) String() string {
	return fmt.Sprintf("BatteryLevel: %.2fV, Temperature: %d°, Type: %s",
		m.BatteryLevel,
		m.Temperature,
		m.FrameType)
}

func calculateBatteryLevel(b byte) float32 {
	return ((3.6-2.8)/255.0)*float32(b) + 2.8
}

func calculateTemperature(b byte) int8 {
	return int8(b)
}

func parseHeader(payload []byte) (*Header, error) {
	return &Header{
		BatteryLevel: calculateBatteryLevel(payload[0]),
		Temperature:  calculateTemperature(payload[1]),
		FrameType:    MoveeType(payload[2])}, nil
}

func parseAlive(payload []byte) (AliveFrame, error) {
	header, err := parseHeader(payload)
	if err != nil {
		return AliveFrame{}, err
	}
	return AliveFrame{Header: header}, nil
}

// c11502151516181717171818181716aa
func parseTemperature(payload []byte) (TemperatureFrame, error) {
	header, err := parseHeader(payload)
	if err != nil {
		return TemperatureFrame{}, err
	}

	var temperatures []int8

	for _, t := range payload[3:] {
		temperatures = append(temperatures, int8(t))
	}

	return TemperatureFrame{Header: header, Temperatures: temperatures}, nil
}

func (t TemperatureFrame) String() string {
	return fmt.Sprintf("%s, Temperatures: %v", t.Header, t.Temperatures)
}

func parseShock(payload []byte) (ShockFrame, error) {
	header, err := parseHeader(payload)
	if err != nil {
		return ShockFrame{}, err
	}
	return ShockFrame{Header: header}, nil
}

// c1150800000082aa or c115080064ffecaa
func parseTilt(payload []byte) (TiltFrame, error) {
	header, err := parseHeader(payload)
	if err != nil {
		return TiltFrame{}, err
	}
	if len(payload) != 7 {
		return TiltFrame{}, fmt.Errorf("Tilt frame expect a payload length of 7, have: %#v", payload)
	}

	var pitchAngle int16
	var rollAngle int16

	if err = binary.Read(bytes.NewReader(payload[3:5]), binary.BigEndian, &pitchAngle); err != nil {
		return TiltFrame{}, err
	}

	if err = binary.Read(bytes.NewReader(payload[5:7]), binary.BigEndian, &rollAngle); err != nil {
		return TiltFrame{}, err
	}

	return TiltFrame{Header: header, PitchAngle: (pitchAngle / 10), RollAngle: (rollAngle / 10)}, nil
}

func (t TiltFrame) String() string {
	return fmt.Sprintf("%s, Pitch angle: %d°, Roll angle: %d°", t.Header, t.PitchAngle, t.RollAngle)
}

// c115100000fffeffe7aa
func parseOrient(payload []byte) (OrientFrame, error) {
	header, err := parseHeader(payload)
	if err != nil {
		return OrientFrame{}, err
	}
	if len(payload) != 9 {
		return OrientFrame{}, fmt.Errorf("Orient frame expect a payload length of 9, have: %#v", payload)
	}

	var pitchAngle int16
	var rollAngle int16
	var yawAngle int16

	if err = binary.Read(bytes.NewReader(payload[3:5]), binary.BigEndian, &pitchAngle); err != nil {
		return OrientFrame{}, err
	}
	if err = binary.Read(bytes.NewReader(payload[5:7]), binary.BigEndian, &rollAngle); err != nil {
		return OrientFrame{}, err
	}
	if err = binary.Read(bytes.NewReader(payload[7:9]), binary.BigEndian, &yawAngle); err != nil {
		return OrientFrame{}, err
	}

	return OrientFrame{Header: header, PitchAngle: pitchAngle, RollAngle: rollAngle, YawAngle: yawAngle}, nil
}

func (o OrientFrame) String() string {
	return fmt.Sprintf("%s, Pitch angle: %d°, Roll angle: %d°, Yaw angle: %d°", o.Header, o.PitchAngle, o.RollAngle, o.YawAngle)
}

// c11a2001aa or c11a2000aa
func parseMotion(payload []byte) (MotionFrame, error) {
	header, err := parseHeader(payload)
	if err != nil {
		return MotionFrame{}, err
	}
	if len(payload) != 4 {
		return MotionFrame{}, fmt.Errorf("Motion frame expect a paylod length of 4, have: %#v", payload)
	}

	onMove := payload[3] == 0x00

	return MotionFrame{Header: header, OnMove: onMove}, nil
}

func (m MotionFrame) String() string {
	return fmt.Sprintf("%s, Motion: %t", m.Header, m.OnMove)
}

// c11a400000000300aa
func parseActivity(payload []byte) (ActivityFrame, error) {
	header, err := parseHeader(payload)
	if err != nil {
		return ActivityFrame{}, err
	}
	if len(payload) != 8 {
		return ActivityFrame{}, fmt.Errorf("Activity frame expect a paylod length of 8, have: %#v", payload)
	}

	onMove := payload[3] == 0x00

	var duration uint32

	if err = binary.Read(bytes.NewReader(payload[4:]), binary.BigEndian, &duration); err != nil {
		return ActivityFrame{}, err
	}

	return ActivityFrame{Header: header, OnMove: onMove, Duration: duration}, nil
}

func (a ActivityFrame) String() string {
	return fmt.Sprintf("%s, Motion: %t, Duration: %dms", a.Header, a.OnMove, a.Duration)
}

// Parse for rotationFrame
func parseRotation(payload []byte) (RotationFrame, error) {
	header, err := parseHeader(payload)
	if err != nil {
		return RotationFrame{}, err
	}
	if len(payload) != 5 {
		return RotationFrame{}, fmt.Errorf("Rotation frame expect a paylod length of 5, have: %#v", payload)
	}

	var nb int16

	if err = binary.Read(bytes.NewReader(payload[3:]), binary.BigEndian, &nb); err != nil {
		return RotationFrame{}, err
	}

	return RotationFrame{Header: header, TurnsNumber: nb}, nil
}

func (r RotationFrame) String() string {
	return fmt.Sprintf("%s, Tour: %d", r.Header, r.TurnsNumber)
}

func parseVibe(payload []byte) (VibeFrame, error) {
	header, err := parseHeader(payload)
	if err != nil {
		return VibeFrame{}, err
	}
	return VibeFrame{Header: header}, nil
}

func parseInformation(payload []byte) (InformationFrame, error) {
	header, err := parseHeader(payload)
	if err != nil {
		return InformationFrame{}, err
	}
	return InformationFrame{Header: header}, nil
}

func parseService(payload []byte) (ServiceFrame, error) {
	header, err := parseHeader(payload)
	if err != nil {
		return ServiceFrame{}, err
	}
	return ServiceFrame{Header: header}, nil
}

// Parse read a slice of data to get the movee frame
func (p Payload) Parse() (MoveeFrame, error) {
	if len(p) < 3 {
		return nil, fmt.Errorf("The payload must have at least 3 bytes: %#v", p)
	}
	switch MoveeType(p[2]) {
	case Alive:
		return parseAlive(p)
	case Temperature:
		return parseTemperature(p)
	case Shock:
		return parseShock(p)
	case Tilt:
		return parseTilt(p)
	case Orient:
		return parseOrient(p)
	case Motion:
		return parseMotion(p)
	case Activity:
		return parseActivity(p)
	case Rotation:
		return parseRotation(p)
	case Vibration:
		return parseVibe(p)
	case Information:
		return parseInformation(p)
	case Service:
		return parseService(p)
	default:
		return nil, fmt.Errorf("Unknown type of frame : %x", p)
	}
}
