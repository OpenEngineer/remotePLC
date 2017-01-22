package blocks

import (
	"../logger/"
	"strconv"
)

// each incoming float64 is converted into a byte (ie. uint8)

type SerialOutput struct {
	OutputBlockData
	address   string
	prevBytes []byte // previous list of bytes sent (compared via string comparison)
}

func (b *SerialOutput) Update() {
	newBytes := FloatsToSerialBytes(b.in)

	// only send a new byte array // TODO: supersede by special node
	if !SerialBytesEqual(newBytes, b.prevBytes) {
		err := SendSerialBytes(b.address, newBytes)
		logger.WriteError("SerialOutput.Update()", err)
		if err == nil {
			//fmt.Println("wrote: ", newBytes)
			b.prevBytes = newBytes
			b.out = b.in
		}
	}
}

func SerialOutputConstructor(name string, words []string) Block {
	address := words[0]
	bitRate, err := strconv.ParseInt(words[1], 10, 64)
	logger.WriteError("SerialOutputConstructor()", err)

	configId := 0
	MakeSerialPort(address, int(bitRate), configId)

	b := &SerialOutput{address: address}
	return b
}

var SerialOutputOk = AddConstructor("SerialOutput", SerialOutputConstructor)
