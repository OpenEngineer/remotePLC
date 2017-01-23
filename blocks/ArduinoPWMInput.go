package blocks

import (
	"../logger/"
	"fmt"
	"strconv"
)

// for more info see the ArduinoPWM.go file
type ArduinoPWMInput struct {
	InputBlockData
	address  string
	question ArduinoPWMPacket
	numBytes int
}

// each byte in the answer packet payload is sent to b.out as a float64
func (b *ArduinoPWMInput) Update() {
	// send the message, and then wait for the reply
	answer, err := SendReceiveArduinoPWM(b.address, b.question)

	fmt.Println("moved past")
	if err == nil {
		bytes := answer.GetPayload()

		if len(bytes) == b.numBytes {
			tmp := SerialBytesToFloats(bytes)

			b.out = tmp
		}
	} else {
		b.out = make([]float64, b.numBytes)
	}

	b.in = b.out
}

func ArduinoPWMInputConstructor(name string, words []string) Block {
	address := words[0]
	bitRate, err := strconv.ParseInt(words[1], 10, 64)

	logger.WriteError("ArduinoPWMInputConstructor()", err)

	numBytes, numBytesErr := strconv.ParseInt(words[2], 10, 64)
	logger.WriteError("ArduinoPWMInputConstructor()", numBytesErr)

	pulseWidth, pulseWidthErr := strconv.ParseInt(words[3], 10, 64)
	logger.WriteError("ArduinoPWMInputConstructor()", pulseWidthErr)

	timeOutCount, timeOutCountErr := strconv.ParseInt(words[4], 10, 64)
	logger.WriteError("ArduinoPWMInputConstructor()", timeOutCountErr)

	// function implemented in ./Serial.go
	configId := 0
	MakeSerialPort(address, int(bitRate), configId)

	b := &ArduinoPWMInput{
		address:  address,
		numBytes: int(numBytes),
		question: ArduinoPWMPacket{
			Header: ArduinoPWMHeader{
				OpCode:       ARDUINOPWM_READOP,
				NumBytes:     uint8(numBytes),
				PulseWidth:   uint16(pulseWidth),
				TimeOutCount: uint16(timeOutCount),
			},
		},
	}

	return b
}

var ArduinoPWMInputOk = AddConstructor("ArduinoPWMInput", ArduinoPWMInputConstructor)
