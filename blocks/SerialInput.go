package blocks

import (
	"../external/mikepb/serial"
	"../logger/"
	"errors"
	"fmt"
	"strconv"
)

type SerialInput struct {
	InputBlockData
	p        *serial.Port
	numInput int
}

func configSerialPortRead(portName string, bitRate int) (*serial.Port, error) {
	options := serial.RawOptions
	options.BitRate = bitRate
	options.Mode = serial.MODE_READ
	options.DataBits = 8
	options.StopBits = 2

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

func (b *SerialInput) Update() {
	nin, err := b.p.InputWaiting()
	if err != nil {
		logger.WriteError("SerialInput.Update()", err)
		return
	} else if nin == 0 {
		return // wait for next update cycle
	}

	buf := make([]byte, nin) // take as many chars as possible from the stream

	n, err := b.p.Read(buf)
	b.p.ResetInput()
	logger.WriteError("SerialInput.Update()", err)

	tmp := make([]float64, b.numInput)
	if n == b.numInput {
		// TODO: use the last packet in the buffer (it is more recent)
		for i, _ := range tmp {
			tmp[i] = float64(buf[i])
		}
		b.out = tmp
	} else {
		b.out = tmp
	}
	b.in = b.out

	fmt.Println(b.out)
}

func SerialInputConstructor(name string, words []string) Block {
	adress := words[0]
	bitRate, err := strconv.ParseInt(words[1], 10, 64)
	logger.WriteError("SerialInputConstructor()", err)

	numInput, numInputErr := strconv.ParseInt(words[2], 10, 64)
	logger.WriteError("SerialInputConstructor()", numInputErr)

	configId := 0
	MakeSerialPort(address, bitRate, configId)

	b := &SerialInput{address: address, numInput: int(numInput)}
	return b
}

var SerialInputOk = AddConstructor("SerialInput", SerialInputConstructor)
