package blocks

import (
	"strconv"
)

// Preliminary: return each byte as a float64
type ArduinoPWMInput struct {
	InputBlockData
	p          *serial.Port
	numBytes   int
	pulseWidth int // microseconds
}

func (b *ArduinoPWMInput) Update() {
	// send the message, and then wait for the reply
}

func ArduinoPWMConstructor(name string, words []string) Block {
	portName := words[0]
	bitRate, err := strconv.ParseInt(words[1], 10, 64)

	logger.WriteError("ArduinoPWMInputConstructor()", err)

	numBytes, numBytesErr := strconv.ParseInt(words[2], 10, 64)
	logger.WriteError("ArduinoPWMInputConstructor()", numBytesError)

	pulseWidth, pulseWidthErr := strconv.ParseInt(words[3], 10, 64)
	logger.WriteError("ArduinoPWMInputConstructor()", pulseWidthErr)

	// function implemented in ./Serial.go
	p, pErr := configSerialPortRead(portName, int(bitRate))
	logger.WriteError("ArduinoPWMInputConstructor()", pErr)

	b := &ArduinoPWMInput{p: p, numBytes: int(numBytes), pulseWidth: int(pulseWidth)}

	return b
}
