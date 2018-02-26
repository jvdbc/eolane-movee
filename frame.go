package frame

import "fmt"

// FrameType data type
type FrameType byte

const (
	// Alive frame
	Alive FrameType = 0x01
	// Temperature frame
	Temperature FrameType = 0x02
	// Shock frame
	Shock FrameType = 0x04
	// Tilt frame
	Tilt FrameType = 0x08
	// Orient frame
	Orient FrameType = 0x10
	// Motion frame
	Motion FrameType = 0x20
	// Activity frame
	Activity FrameType = 0x40
	// Rotation frame
	Rotation FrameType = 0x80
	// Vibration frame
	Vibration FrameType = 0x86
	// Information frame
	Information FrameType = 0xFE
	// Service frame
	Service FrameType = 0xFF
)

// Payload type alias on slice of bytes
type Payload []byte

// Header type
type Header struct {
	BatteryLevel float32
	Temperature  int8
	FrameType    FrameType
}

// MoveeFrame interface
type MoveeFrame interface {
	GetHeader() *Header
}

func parseHeader(payload []byte) (*Header, error) {
	if len(payload) < 3 {
		return &Header{}, fmt.Errorf("The payload must have at least 3 bytes: %x", payload)
	}

	return &Header{
		BatteryLevel: calculateBatteryLevel(payload[0]),
		Temperature:  calculateTemperature(payload[1]),
		FrameType:    FrameType(payload[2])}, nil
}

// AliveFrame type
type AliveFrame struct {
	*Header
}

// GetHeader for AliveFrame
func (f AliveFrame) GetHeader() *Header {
	return f.Header
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

// GetHeader for TemperatureFrame
func (f TemperatureFrame) GetHeader() *Header {
	return f.Header
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

// GetHeader for ShockFrame
func (f ShockFrame) GetHeader() *Header {
	return f.Header
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

// GetHeader for TiltFrame
func (f TiltFrame) GetHeader() *Header {
	return f.Header
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

// GetHeader for OrientFrame
func (f OrientFrame) GetHeader() *Header {
	return f.Header
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

// GetHeader for OrientFrMotionFrameame
func (f MotionFrame) GetHeader() *Header {
	return f.Header
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

// GetHeader for ActivityFrame
func (f ActivityFrame) GetHeader() *Header {
	return f.Header
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

// GetHeader for RotationFrame
func (f RotationFrame) GetHeader() *Header {
	return f.Header
}

// Parse for rotationFrame
func parseRotation(payload []byte) (RotationFrame, error) {
	header, err := parseHeader(payload)

	if err != nil {
		return RotationFrame{}, err
	}

	if len(payload) != 4 {
		return RotationFrame{}, fmt.Errorf("Rotation frame has a expected length of 4: %x", payload)
	}

	return RotationFrame{Header: header, numberOfTurns: int16(payload[3])}, nil
}

// VibeFrame type
type VibeFrame struct {
	*Header
}

// GetHeader for VibeFrame
func (f VibeFrame) GetHeader() *Header {
	return f.Header
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

// GetHeader for InformationFrame
func (f InformationFrame) GetHeader() *Header {
	return f.Header
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

// GetHeader for ServiceFrame
func (f ServiceFrame) GetHeader() *Header {
	return f.Header
}

func parseService(payload []byte) (ServiceFrame, error) {
	header, err := parseHeader(payload)

	if err != nil {
		return ServiceFrame{}, err
	}

	return ServiceFrame{Header: header}, nil
}

func calculateBatteryLevel(b byte) float32 {
	return ((3.6-2.8)/255.0)*float32(b) + 2.8
}

func calculateTemperature(b byte) int8 {
	return int8(b)
}

// Parse read a slice of data to get the movee frame
func (p Payload) Parse() (MoveeFrame, error) {
	if p == nil {
		return nil, fmt.Errorf("The payload must not be nil")
	}

	if len(p) < 3 {
		return nil, fmt.Errorf("The payload must have at least 3 bytes: %x", p)
	}

	switch FrameType(p[2]) {
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
