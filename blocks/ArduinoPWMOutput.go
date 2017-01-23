package blocks

import (
	"../logger/"
	//"fmt"
	"strconv"
)

// for more info see the ArduinoPWM.go file
type ArduinoPWMOutput struct {
	OutputBlockData
	address  string
	question ArduinoPWMPacket
}

func (b *ArduinoPWMOutput) Update() {
	numBytes := len(b.in)
	b.question.Header.NumBytes = uint8(numBytes)

	for i, v := range b.in {
		b.question.Bytes[i] = byte(uint8(v))
	}

	//fmt.Println("sending message")
	_, err := SendReceiveArduinoPWM(b.address, b.question) // ORIG
	//SendArduinoPWM(b.address, b.question)

	if err == nil { // ORIG
		b.out = b.in
	} // ORIG
}

func ArduinoPWMOutputConstructor(name string, words []string) Block {
	address := words[0]
	bitRate, err := strconv.ParseInt(words[1], 10, 64)
	logger.WriteError("ArduinoPWMOutputConstructor()", err)

	// the pulsewidth is specific to a set of eg. 433MHz devices
	pulseWidth, pulseWidthErr := strconv.ParseInt(words[2], 10, 64)
	logger.WriteError("ArduinoPWMOutputConstructor()", pulseWidthErr)

	// function implemented in ./Serial.go
	configId := 0
	MakeSerialPort(address, int(bitRate), configId)

	b := &ArduinoPWMOutput{
		address: address,
		question: ArduinoPWMPacket{
			Header: ArduinoPWMHeader{
				OpCode:     ARDUINOPWM_WRITEOP,
				PulseWidth: uint16(pulseWidth),
			},
		},
	}

	return b
}

var ArduinoPWMOutputOk = AddConstructor("ArduinoPWMOutput", ArduinoPWMOutputConstructor)
