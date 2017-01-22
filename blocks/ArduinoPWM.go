package blocks

import (
	"../logger/"
	"bytes"
	"encoding/binary"
	"errors"
)

// you must compile the remoteEmbeddedSystems/arduino/duplexPWM/arduinoDuplexPWM.cpp file
// and install it on the arduino

const (
	ARDUINOPWM_MAX_BYTES int   = 255 // largest numBytes is header
	ARDUINOPWM_WRITEOP   uint8 = 1
	ARDUINOPWM_READOP    uint8 = 2
)

type ArduinoPWMHeader struct {
	OpCode       uint8  // eg WRITE_OPCODE or READ_OPCODE
	NumBytes     uint8  // size of payload
	PulseWidth   uint16 // duration in microseconds of smallest single pulse
	TimeOutCount uint16 // only for READ_OPCODE, stop trying to read a message after this number of pulses
	ErrorCode    uint8  // returned by arduino, set to 0 when sending message to arduino
}

type ArduinoPWMPacket struct {
	Header ArduinoPWMHeader
	Bytes  [ARDUINOPWM_MAX_BYTES]byte
}

// includes header
func (p *ArduinoPWMPacket) Size() int {
	numBytes := 7 + int(p.Header.NumBytes)
	return numBytes
}

func (p *ArduinoPWMPacket) GetPayload() []byte {
	numBytes := p.Header.NumBytes
	return p.Bytes[0:numBytes]
}

func arduinoPWMPacketToBytes(p ArduinoPWMPacket) []byte {
	b := new(bytes.Buffer)
	err := binary.Write(b, binary.LittleEndian, p)
	logger.WriteError("arduinoPWMPacketToBytes()", err)

	numBytes := p.Size()

	return b.Bytes()[0:numBytes]
}

func arduinoPWMBytesToPacket(b []byte) ArduinoPWMPacket {
	var p ArduinoPWMPacket
	buffer := bytes.NewBuffer(b)
	binary.Read(buffer, binary.LittleEndian, &p)

	return p
}

// shoot and forget message
func SendArduinoPWM(address string, p ArduinoPWMPacket) {
	b := arduinoPWMPacketToBytes(p)

	SendSerialBytes(address, b)
}

func SendReceiveArduinoPWM(address string, p0 ArduinoPWMPacket) (ArduinoPWMPacket, error) {
	b0 := arduinoPWMPacketToBytes(p0)

	b1, err := SendReceiveSerialBytes(address, b0, p0.Size())

	logger.WriteError("SendReceiveArduinoPWM()", err)

	p1 := arduinoPWMBytesToPacket(b1)

	if p1.Header.ErrorCode != 0 {
		err = errors.New("nonzero error code from arduino")
		logger.WriteError("SendReceiveArduinoPWM()", err)
	}

	return p1, err
}
