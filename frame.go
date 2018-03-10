package frame

import "fmt"

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

func (m Header) String() string {
	return fmt.Sprintf("BatteryLevel: %f, Temperature: %d, Type: %s",
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

// AliveFrame type
type AliveFrame struct {
	*Header
}

func parseAlive(payload []byte) (AliveFrame, error) {
	header, err := parseHeader(payload)
	if err != nil {
		return AliveFrame{}, err
	}
	return AliveFrame{Header: header}, nil
}

// TemperatureFrame type
type TemperatureFrame struct {
	*Header
}

func parseTemperature(payload []byte) (TemperatureFrame, error) {
	header, err := parseHeader(payload)
	if err != nil {
		return TemperatureFrame{}, err
	}
	return TemperatureFrame{Header: header}, nil
}

// ShockFrame type
type ShockFrame struct {
	*Header
}

func parseShock(payload []byte) (ShockFrame, error) {
	header, err := parseHeader(payload)
	if err != nil {
		return ShockFrame{}, err
	}
	return ShockFrame{Header: header}, nil
}

// TiltFrame type
// c1150800000000aac115080000ffecaac115080064ffecaa
type TiltFrame struct {
	*Header
}

func parseTilt(payload []byte) (TiltFrame, error) {
	header, err := parseHeader(payload)
	if err != nil {
		return TiltFrame{}, err
	}
	return TiltFrame{Header: header}, nil
}

// OrientFrame type
type OrientFrame struct {
	*Header
}

func parseOrient(payload []byte) (OrientFrame, error) {
	header, err := parseHeader(payload)
	if err != nil {
		return OrientFrame{}, err
	}
	return OrientFrame{Header: header}, nil
}

// MotionFrame type
type MotionFrame struct {
	*Header
}

func parseMotion(payload []byte) (MotionFrame, error) {
	header, err := parseHeader(payload)
	if err != nil {
		return MotionFrame{}, err
	}
	return MotionFrame{Header: header}, nil
}

// ActivityFrame type
type ActivityFrame struct {
	*Header
}

func parseActivity(payload []byte) (ActivityFrame, error) {
	header, err := parseHeader(payload)
	if err != nil {
		return ActivityFrame{}, err
	}
	return ActivityFrame{Header: header}, nil
}

// RotationFrame type
// c11580000aaa
type RotationFrame struct {
	*Header
	numberOfTurns int16
}

// Parse for rotationFrame
func parseRotation(payload []byte) (RotationFrame, error) {
	header, err := parseHeader(payload)
	if err != nil {
		return RotationFrame{}, err
	}
	if len(payload) != 5 {
		return RotationFrame{}, fmt.Errorf("Rotation frame has a expected length of 5: %#v", payload)
	}
	return RotationFrame{Header: header, numberOfTurns: int16(payload[4])}, nil
}

// VibeFrame type
type VibeFrame struct {
	*Header
}

func parseVibe(payload []byte) (VibeFrame, error) {
	header, err := parseHeader(payload)
	if err != nil {
		return VibeFrame{}, err
	}
	return VibeFrame{Header: header}, nil
}

// InformationFrame type
type InformationFrame struct {
	*Header
}

func parseInformation(payload []byte) (InformationFrame, error) {
	header, err := parseHeader(payload)
	if err != nil {
		return InformationFrame{}, err
	}
	return InformationFrame{Header: header}, nil
}

// ServiceFrame type
type ServiceFrame struct {
	*Header
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
