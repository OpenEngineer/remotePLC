package blocks

import (
	"../logger/"
	"../thirdParty/mikepb/serial"
	"errors"
	//"fmt"
	"strconv"
)

// each incoming float64 is converted into a byte (ie. uint8)

type SerialOutput struct {
	OutputBlockData
	p *serial.Port
	//portName string
	//baudRate  int
	prevBytes []byte // previous list of bytes sent (compared via string comparisson)
}

func floats2Bytes(x []float64) []byte {
	b := make([]byte, len(x))

	for i, v := range x {
		b[i] = byte(uint8(v))
	}

	return b
}

func bytesEqual(b1 []byte, b2 []byte) bool {
	isEqual := true
	if len(b1) == len(b2) {
		for i, v := range b1 {
			if uint8(v) != uint8(b2[i]) {
				isEqual = false
			}
		}
	} else {
		isEqual = false
	}

	return isEqual

}

func listSerialPorts() {
	info, err := serial.ListPorts()
	logger.WriteError("listSerialPorts()", err)

	for _, v := range info {
		logger.WriteEvent("possible serial port: ", v.Name())
	}
}

func configSerialPort(portName string, bitRate int) (*serial.Port, error) {
	options := serial.RawOptions
	options.BitRate = bitRate
	options.Mode = serial.MODE_WRITE
	options.DataBits = 8
	options.StopBits = 1

	p, openErr := options.Open(portName)
	if openErr != nil {
		logger.WriteEvent("problem opening " + portName + " in configSerialPort(), you might have insufficient rights")
		return nil, openErr
	}

	if p == nil {
		logger.WriteEvent("actual serial port: ", portName)
		listSerialPorts()
		logger.WriteFatal("configSerialPort()", errors.New("nil port"))
	}

	applyErr := p.Apply(&options)
	if applyErr != nil {
		return nil, applyErr
	}

	return p, nil
}

func (b *SerialOutput) Update() {
	newBytes := floats2Bytes(b.in)

	// only send a new byte array
	if !bytesEqual(newBytes, b.prevBytes) {

		_, writeErr := b.p.Write(newBytes)
		logger.WriteError("SerialOutput.Update()", writeErr)
		if writeErr == nil {
			//fmt.Println("wrote: ", newBytes)
			b.prevBytes = newBytes
			b.out = b.in
		}
	}
}

func SerialOutputConstructor(name string, words []string) Block {
	portName := words[0]
	bitRate, err := strconv.ParseInt(words[1], 10, 64)

	logger.WriteError("SerialOutputConstructor()", err)

	p, pErr := configSerialPort(portName, int(bitRate))
	logger.WriteError("SerialOutputConstructor()", pErr)

	b := &SerialOutput{p: p}
	return b
}

var SerialOutputOk = AddConstructor("SerialOutput", SerialOutputConstructor)
