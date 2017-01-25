package blocks

import (
	"../logger/"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"sync"
	"time"
)

// you must compile the remoteEmbeddedSystems/arduino/duplexPWM/arduinoDuplexPWM.cpp file
// and install it on the arduino

const (
	ARDUINOPWM_MAX_BYTES int   = 255 // largest numBytes in header
	ARDUINOPWM_WRITEOP   uint8 = 1
	ARDUINOPWM_READOP    uint8 = 2
)

type ArduinoPWMHeader struct {
	OpCode       uint8  // eg WRITE_OPCODE or READ_OPCODE
	NumBytes     uint8  // size of payload
	PulseWidth   uint16 // duration in microseconds of smallest single pulse
	ClearCount   uint8  // number of pulses the line should be clear before recording, only for READ_OPCODE (ignored otherwise)
	TimeOutCount uint16 // only for READ_OPCODE (ignored otherwise), stop trying to read a message after this number of pulses
	ErrorCode    uint8  // returned by arduino, set to 0 when sending message to arduino
}

type ArduinoPWMPacket struct {
	Header ArduinoPWMHeader
	Bytes  [ARDUINOPWM_MAX_BYTES]byte
}

// includes header
func (p *ArduinoPWMPacket) Size() int {
	numBytes := 8 + int(p.Header.NumBytes)
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
func sendArduinoPWM(address string, p ArduinoPWMPacket) {
	b := arduinoPWMPacketToBytes(p)

	SendSerialBytes(address, b)
}

// these single instance buffers should be locked before being written to
var arduinoPWMPacketBuffer ArduinoPWMPacket
var arduinoPWMPacketBufferMutex sync.Mutex

var arduinoPWMSendingMutex sync.Mutex
var arduinoPWMSendingState bool
var arduinoPWMSendingStateMutex sync.Mutex

func setArduinoPWMPacketBuffer(p ArduinoPWMPacket) {
	arduinoPWMPacketBufferMutex.Lock()
	arduinoPWMPacketBuffer = p
	arduinoPWMPacketBufferMutex.Unlock()
}

func resetArduinoPWMPacketBuffer() {
	setArduinoPWMPacketBuffer(ArduinoPWMPacket{})
}

func getArduinoPWMPacketBuffer() ArduinoPWMPacket {
	arduinoPWMPacketBufferMutex.Lock()
	p := arduinoPWMPacketBuffer
	arduinoPWMPacketBufferMutex.Unlock()

	return p
}

func lockArduinoPWMSending() {
	arduinoPWMSendingMutex.Lock()
	arduinoPWMSendingStateMutex.Lock()
	arduinoPWMSendingState = true
	arduinoPWMSendingStateMutex.Unlock()
}

func unlockArduinoPWMSending() {
	arduinoPWMSendingStateMutex.Lock()
	arduinoPWMSendingState = false
	arduinoPWMSendingStateMutex.Unlock()
	arduinoPWMSendingMutex.Unlock()
}

func lockIfUnlockedArduinoPWMSending() bool {
	ok := false

	arduinoPWMSendingStateMutex.Lock()

	if !arduinoPWMSendingState {
		arduinoPWMSendingMutex.Lock()
		arduinoPWMSendingState = true
		ok = true
	}

	arduinoPWMSendingStateMutex.Unlock()

	return ok
}

// TODO: take timeouts into account
func sendReceiveArduinoPWMPacket(address string, p0 ArduinoPWMPacket) (ArduinoPWMPacket, error) {
	b0 := arduinoPWMPacketToBytes(p0)

	// timeout after pulseWidth*timeOutCount
	timeOutDuration := time.Duration(int(p0.Header.PulseWidth) *
		int(p0.Header.TimeOutCount) * 1000)
	b1, err := SendReceiveSerialBytes(address, b0, p0.Size(), time.Now().Add(timeOutDuration))

	fmt.Println("received ", b1)

	if err != nil {
		logger.WriteEvent("SendReceiveArduinoPWM()", err)
	}

	p1 := arduinoPWMBytesToPacket(b1)

	if p1.Header.ErrorCode != 0 {
		err = errors.New("nonzero error code from arduino")
		logger.WriteError("SendReceiveArduinoPWM()", err)
	}

	return p1, err
}

// TODO: return response packet once this functionality is needed
// for now: shoot and forget message
func sendReceiveArduinoPWMWriteOp(address string, p0 ArduinoPWMPacket) {
	if ok := lockIfUnlockedArduinoPWMSending(); ok {
		sendArduinoPWM(address, p0)
		resetArduinoPWMPacketBuffer()
		unlockArduinoPWMSending()
	} else {
		setArduinoPWMPacketBuffer(p0)
	}
}

func sendReceiveArduinoPWMReadOp(address string, p0 ArduinoPWMPacket) (ArduinoPWMPacket, error) {
	lockArduinoPWMSending()

	p1, err := sendReceiveArduinoPWMPacket(address, p0)

	pw := getArduinoPWMPacketBuffer()

	if pw.Header.OpCode == ARDUINOPWM_WRITEOP {
		sendArduinoPWM(address, pw)
		resetArduinoPWMPacketBuffer()
	}

	unlockArduinoPWMSending()

	return p1, err
}

// depending on the packet opCode, do something different
func SendReceiveArduinoPWM(address string, p0 ArduinoPWMPacket) (ArduinoPWMPacket, error) {
	var p1 ArduinoPWMPacket
	var err error

	switch op := p0.Header.OpCode; op {
	case ARDUINOPWM_WRITEOP:
		sendReceiveArduinoPWMWriteOp(address, p0)
	case ARDUINOPWM_READOP:
		p1, err = sendReceiveArduinoPWMReadOp(address, p0)
	default:
		return ArduinoPWMPacket{}, errors.New("opCode not recognized")
	}

	return p1, err

}
