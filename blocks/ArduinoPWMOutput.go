package blocks

import (
	"../logger/"
	"strconv"
)

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

	_, err := SendReceiveArduinoPWM(b.address, b.question)

	if err == nil {
		b.out = b.in
	}
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
