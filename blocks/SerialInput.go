package blocks

import (
	"../logger/"
	//"fmt"
	"strconv"
)

type SerialInput struct {
	InputBlockData
	address  string
	numInput int
}

func (b *SerialInput) Update() {
	bytes, _ := ReceiveAllSerialBytes(b.address)
	n := len(bytes)

	// convert to floats
	tmp := make([]float64, b.numInput)
	if n == b.numInput {
		for i, _ := range tmp {
			tmp[i] = float64(bytes[n-b.numInput+i]) // use last packet in buffer
		}
		b.out = tmp
	} else { // return an array with zeros
		b.out = tmp
	}
	b.in = b.out

	//fmt.Println(b.out)
}

func SerialInputConstructor(name string, words []string) Block {
	address := words[0]
	bitRate, err := strconv.ParseInt(words[1], 10, 64)
	logger.WriteError("SerialInputConstructor()", err)

	numInput, numInputErr := strconv.ParseInt(words[2], 10, 64)
	logger.WriteError("SerialInputConstructor()", numInputErr)

	configId := 0
	MakeSerialPort(address, int(bitRate), configId)

	b := &SerialInput{address: address, numInput: int(numInput)}
	return b
}

var SerialInputOk = AddConstructor("SerialInput", SerialInputConstructor)
