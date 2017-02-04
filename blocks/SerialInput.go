package blocks

import (
	"../logger/"
	//"fmt"
	"../parser/"
	"time"
)

type SerialInput struct {
	InputBlockData
	address  string
	numInput int
}

func (b *SerialInput) Poll() {
	bytes, _ := ReceiveAllSerialBytes(b.address) // assumed to have some delay
	n := len(bytes)

	// convert to floats
	b.in = MakeUndefined(b.numInput)
	if n == b.numInput {
		for i, _ := range b.in {
			b.in[i] = float64(bytes[n-b.numInput+i]) // use last packet in buffer
		}
	}
}

func (b *SerialInput) Update() {
	if len(b.in) == b.numInput {
		b.out = SafeCopy(b.numInput, b.in, b.numInput)
	} else {
		b.out = MakeUndefined(b.numInput)
	}

	b.in = MakeUndefined(b.numInput) // only a new poll reply gives us defined numbers
}

// TODO: add delay
func SerialInputConstructor(name string, words []string) Block {
	var address string
	var bitRate int
	var numInput int
	var periodString string

	positional := parser.PositionalArgs(&address, &bitRate, &numInput, &periodString)
	parser.ParsePositionalArgs(words, positional)

	configId := 0
	MakeSerialPort(address, bitRate, configId)

	b := &SerialInput{address: address, numInput: numInput}

	// poll infinitely in the background
	period, err := time.ParseDuration(periodString)
	logger.WriteError("SerialInputConstructor()", err)
	go func() {
		ticker := time.NewTicker(period)
		for {
			<-ticker.C
			b.Poll()
		}
	}()

	return b
}

var SerialInputOk = AddConstructor("SerialInput", SerialInputConstructor)
