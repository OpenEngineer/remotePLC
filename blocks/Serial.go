package blocks

// contains all general Serial connection structures and variables

import (
	"../external/serial"
	"../logger/"
	"errors"
	//"fmt"
	"time"
)

// wrapper around serial.Port
// contains construction data as well, for easy access
type SerialPort struct {
	address  string
	bitRate  int
	configId int // within RS232 many configuration are possible (eg. number of stop bits, parity bit ...)

	p *serial.Port
}

// each constructed port has a unique address string
var serialPorts = make(map[string]SerialPort)

// different configurations are possible,
//  these are loaded dynamically
//  bitRates are kept separate because I expect them to require more frequent changes
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
			logger.WriteFatal("MakeSerialPort()",
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

func GetSerialPort(address string) (SerialPort, error) {
	var err error
	var p SerialPort
	var ok bool

	if p, ok = serialPorts[address]; ok {
		err = nil
	} else {
		err = errors.New("port " + address + " not found, ignoring")
	}

	return p, err
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

func SerialBytesToFloats(b []byte) []float64 {
	f := make([]float64, len(b))

	for i, v := range b {
		f[i] = float64(uint8(v))
	}

	return f
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
	p, err := GetSerialPort(address)
	if err == nil {
		_, err = p.p.Write(bytes)
	}

	return err
}

// receive message of numBytes size
// smaller messages are ignored and left in the serialPort buffer
func ReceiveSerialBytes(address string, numBytes int) ([]byte, error) {
	p, _ := GetSerialPort(address)

	buffer := make([]byte, numBytes)

	_, err := p.p.Read(buffer) // read the message of fixed size

	logger.WriteError("ReceiveSerialBytes()", err)

	return buffer, err
}

func ReceiveAllSerialBytes(address string) ([]byte, error) {
	p, _ := GetSerialPort(address)

	numIncoming, err := p.p.InputWaiting()
	if err != nil {
		logger.WriteError("ReceiveAllSerialBytes()", err)
		return []byte{}, err
	} else if numIncoming == 0 {
		return []byte{}, nil
	}

	buffer := make([]byte, numIncoming) // take as many chars as possible from the stream

	numRead, err := p.p.Read(buffer) // read everything
	p.p.ResetInput()                 // throw away everything in the serialPort buffer

	logger.WriteError("ReceiveAllSerialBytes()", err)

	return buffer[0:numRead], err
}

// send and wait for a reply of fixed length
// this function can be used for synchronous communication
func SendReceiveSerialBytes(address string, bytes []byte, numBytes int, deadline time.Time) ([]byte, error) {
	p, err := GetSerialPort(address)

	// set the deadline
	p.p.SetDeadline(deadline)

	// the reply will go into this buffer:
	buffer := make([]byte, numBytes)

	// if the port is ok, write and read the messages
	if err == nil {
		var errWrite error

		p.p.Reset()

		_, errWrite = p.p.Write(bytes)

		//fmt.Println("sent     ", bytes)
		if errWrite != nil {
			logger.WriteError("SendReceiveSerialBytes(), write", errWrite)
			return []byte{}, errWrite
		}

		p.p.Sync()

		_, errRead := p.p.Read(buffer)
		if errRead != nil {
			logger.WriteError("SendReceiveSerialBytes(), read", errRead)
			return []byte{}, errRead
		}

		return buffer, nil
	} else {
		return []byte{}, err
	}

	// unset the deadline
	p.p.SetDeadline(time.Time{})

	return buffer, nil
}
