package blocks

import (
	"../logger/"
	"../thirdParty/mikepb/serial"
	"strconv"
)

// each incoming float64 is converted into a byte (ie. uint8)

type SerialOutput struct {
	OutputBlockData
	portName  string
	baudRate  int
	prevBytes []byte // previous list of bytes sent (compared via string comparisson)
}

func floats2Bytes(x []float64) []byte {
	b := make([]byte, len(x))

	for i, v := range x {
		b[i] = byte(uint8(v))
	}

	return b
}

func (b *SerialOutput) configPort() (*serial.Port, error) {
	p, openErr := serial.Open(b.portName)
	if openErr != nil {
		return nil, openErr
	}

	options := serial.Options{
		Mode:        serial.MODE_WRITE,
		BitRate:     b.baudRate,
		DataBits:    8,
		StopBits:    1,
		Parity:      serial.PARITY_NONE,
		FlowControl: serial.FLOWCONTROL_NONE,
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
	if string(newBytes) != string(b.prevBytes) {
		p, openErr := b.configPort()
		if openErr != nil {

			_, writeErr := p.Write(newBytes)
			logger.WriteError("SerialOutput.Update()", writeErr)
			if writeErr != nil {

				b.prevBytes = newBytes
				b.out = b.in
			}
		}
	}
}

func SerialOutputConstructor(name string, words []string) Block {
	portName := words[0]
	baudRate, err := strconv.ParseInt(words[1], 10, 64)

	logger.WriteError("SerialOutputConstructor()", err)

	b := &SerialOutput{portName: portName, baudRate: int(baudRate)}
	return b
}

var SerialOutputOk = AddConstructor("SerialOutput", SerialOutputConstructor)
