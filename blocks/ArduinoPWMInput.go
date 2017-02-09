package blocks

import (
	//"../logger/"
	"../parser/"
	//"strconv"
	//"fmt"
	"time"
)

// for more info see the ArduinoPWM.go file
type ArduinoPWMInput struct {
	InputBlockData
	address  string
	question ArduinoPWMPacket
	numInput int
}

// each byte in the answer packet payload is sent to b.out as a float64
func (b *ArduinoPWMInput) Poll() {
	// send the message, and then wait for the reply
	answer, err := SendReceiveArduinoPWM(b.address, b.question) // this function handles interlacing of read and write requests from different blocks

	//fmt.Println(answer.GetPayload(), b.question)
	if err == nil {
		bytes := answer.GetPayload()

		if len(bytes) == b.numInput {
			tmp := SerialBytesToFloats(bytes)

			b.in = tmp
		} else {
			b.in = MakeUndefined(b.numInput)
		}
	} else {
		b.in = MakeUndefined(b.numInput)
	}
}

func (b *ArduinoPWMInput) Update() {
	if len(b.in) == b.numInput {
		b.out = SafeCopy(b.numInput, b.in, b.numInput)
	} else {
		b.out = MakeUndefined(b.numInput)
	}

	b.in = MakeUndefined(b.numInput) // only a new poll reply gives us defined numbers
}

func ArduinoPWMInputConstructor(name string, words []string) Block {
	var address string
	var bitRate int
	var numInput int
	var pulseWidth int
	var clearCount int
	var timeOutCount int
	var pulseMargin int

	positional := parser.PositionalArgs(&address, &bitRate, &numInput, &pulseWidth, &clearCount, &timeOutCount, &pulseMargin)
	optional := parser.OptionalArgs()

	parser.ParseArgs(words, positional, optional)

	// function implemented in ./Serial.go
	configId := 0
	MakeSerialPort(address, bitRate, configId)

	b := &ArduinoPWMInput{
		address:  address,
		numInput: numInput,
		question: ArduinoPWMPacket{
			Header1: ArduinoPWMHeader1{
				OpCode:     ARDUINOPWM_READOP,
				NumBytes:   uint8(numInput),
				PulseWidth: uint16(pulseWidth),
			},
			Header2: ArduinoPWMHeader2{
				ClearCount:   uint8(clearCount),
				TimeOutCount: uint16(timeOutCount),
				PulseMargin:  uint8(pulseMargin),
			},
		},
	}

	// loop the polling infinitely in the background
	go func() {
		period, _ := time.ParseDuration("1s")
		ticker := time.NewTicker(period)
		for {
			<-ticker.C
			b.Poll()
		}
	}()

	return b
}

var ArduinoPWMInputOk = AddConstructor("ArduinoPWMInput", ArduinoPWMInputConstructor)
