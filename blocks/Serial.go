package blocks

// contains all general Serial connection structures and variables

import (
	"../external/mikepb/serial"
	"../logger/"
	"errors"
)

// wrapper around serial.Port
// contains construction data as well, for easy access
type SerialPort struct {
	address  string
	bitRate  int
	configId int // protocolId itself is RS232 (at least how I think of it)

	p *serial.Port
}

// with this map we make sure that each port has a unique address
var serialPorts = make(map[string]SerialPort)

// different configurations are possible,
//  these are loaded dynamically
//  bitRates are kept separate because except them to require more frequent changes
var serialPortConfigs = make(map[int]func(string, int) (*serial.Port, error))

func (p *SerialPort) isEqual(bitRate int, configId int) bool {
	if p.bitRate == bitRate && p.configId == configId {
		return true
	} else {
		return false
	}
}

func AddSerialPortConfig(configId int, fn func(string, int) (*serial.Port, error)) bool {
	serialPortConfigs[configId] = fn
	return true
}

func listSerialPorts() {
	info, err := serial.ListPorts()
	logger.WriteError("listSerialPorts()", err)

	for _, v := range info {
		logger.WriteEvent("possible serial port: ", v.Name())
	}
}

func defaultSerialConfig(address string, bitRate int) (*serial.Port, error) {
	options := serial.RawOptions
	options.BitRate = bitRate
	options.Mode = serial.MODE_READ_WRITE
	options.DataBits = 8
	options.StopBits = 2

	p, openErr := options.Open(address)
	if openErr != nil {
		logger.WriteEvent("problem opening " + address + " in defaultSerialConfig(), you might have insufficient rights")
		return nil, openErr
	}

	if p == nil {
		logger.WriteEvent("actual serial port: ", address)
		listSerialPorts()
		logger.WriteFatal("defaultSerialConfig()", errors.New("nil port"))
	}

	applyErr := p.Apply(&options)
	if applyErr != nil {
		return nil, applyErr
	}

	return p, nil
}

var defaultSerialConfigOk = AddSerialPortConfig(0, defaultSerialConfig)

// example address="/dev/ACM0tty"
// example bitRate=9600
// configId: 0 default (2 stops bits, etc)
//                no others supported yet
//  (I just included this to warn that different protocol variants are possible)
func MakeSerialPort(address string, bitRate int, configId int) {
	// has a port at this address already been opened?
	if p, ok := serialPorts[address]; ok {
		if !p.isEqual(bitRate, configId) {
			log.WriteFatal("MakeSerialPort()",
				errors.New("same address serialPort for different bitRate or configId not allowed"))
		}
	} else { // configure a new port
		p, err := serialPortConfigs[configId](address, bitRate)

		logger.WriteError("MakeSerialPort()", err)

		serialPorts[address] = SerialPort{
			address:  address,
			bitRate:  bitRate,
			configId: configId,
			p:        p,
		}
	}
}

// functions for sending and receiving bytes from serial port

// some message creation and comparison functions
func FloatsToSerialBytes(x []float64) []byte {
	b := make([]byte, len(x))

	for i, v := range x {
		b[i] = byte(uint8(v))
	}

	return b
}

func SerialBytesEqual(b1 []byte, b2 []byte) bool {
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

// the main send and receive functions

// send without caring about reply message
func SendSerialBytes(address string, bytes []byte) error {
	var err error
	if p, ok := serialPorts[address]; ok {
		_, err = p.p.Write(bytes)
	} else {
		err = errors.New("port ", address, " not found, ignoring")
	}

	return err
}

// receive anything
func ReceiveSerialBytes(address string, numBytes int, bytes []byte) (int, error) {
}
