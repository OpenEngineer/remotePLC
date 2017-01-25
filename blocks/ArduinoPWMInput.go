package blocks

import (
	//"../logger/"
	"../parser/"
	"fmt"
	//"strconv"
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
	var address string
	var bitRate int
	var numBytes int
	var pulseWidth int
	var clearCount int
	var timeOutCount int

	positional := parser.PositionalArgs(&address, &bitRate, &numBytes, &pulseWidth, &clearCount, &timeOutCount)
	optional := parser.OptionalArgs()

	parser.ParseArgs(words, positional, optional)

	// function implemented in ./Serial.go
	configId := 0
	MakeSerialPort(address, bitRate, configId)

	b := &ArduinoPWMInput{
		address:  address,
		numBytes: numBytes,
		question: ArduinoPWMPacket{
			Header: ArduinoPWMHeader{
				OpCode:       ARDUINOPWM_READOP,
				NumBytes:     uint8(numBytes),
				PulseWidth:   uint16(pulseWidth),
				ClearCount:   uint8(clearCount),
				TimeOutCount: uint16(timeOutCount),
			},
		},
	}

	return b
}

var ArduinoPWMInputOk = AddConstructor("ArduinoPWMInput", ArduinoPWMInputConstructor)
