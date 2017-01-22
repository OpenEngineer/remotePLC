package blocks

// code used by both ExampleUDPInput.go and ExampleUDPOutput.go

import (
	"../logger/"
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
)

const (
	EXAMPLEUDP_MAX_BYTES   int    = 1460 // the protocol must make sure that the EXAMPLEUDP_MAX_RECORDS can always fit in this
	EXAMPLEUDP_MAX_RECORDS int    = 10
	EXAMPLEUDP_DATA_SIZE   int    = 4
	EXAMPLEUDP_PORT        string = "666"
	EXAMPLEUDP_VERSION     byte   = 0x01
	EXAMPLEUDP_READOP      byte   = 0x01
	EXAMPLEUDP_WRITEOP     byte   = 0x02
)

type ExampleUDPHeader struct {
	Version    byte
	OpId       byte // 0x01 for read, 0x02 for write
	SeqNo      uint16
	NumRecords uint16
	ErrorCode  uint16 // 0 in case of no error, not used when sending a question
}

type ExampleUDPRecord struct {
	Id   uint16                     // example: id/address on target machine
	Data [EXAMPLEUDP_DATA_SIZE]byte // example: a float32
}

type ExampleUDPPacket struct {
	Header  ExampleUDPHeader
	Records [EXAMPLEUDP_MAX_RECORDS]ExampleUDPRecord
}

// includes the size of the header
func (p *ExampleUDPPacket) Size() int {
	numRecords := int(p.Header.NumRecords)
	numBytes := 8 + numRecords*(EXAMPLEUDP_DATA_SIZE+2)
	return numBytes
}

func exampleUDPPacketToBytes(p ExampleUDPPacket) []byte {
	b := new(bytes.Buffer)
	err := binary.Write(b, binary.LittleEndian, p)
	logger.WriteError("exampleUDPPacketToBytes()", err)

	numBytes := p.Size()

	return b.Bytes()[0:numBytes]
}

func exampleUDPBytesToPacket(b []byte) ExampleUDPPacket {
	var p ExampleUDPPacket
	buffer := bytes.NewBuffer(b)
	binary.Read(buffer, binary.LittleEndian, &p)

	return p
}

var exampleUDPConnection net.Conn // shared by ExampleUDPInput and ExampleUDPOutput

// a common port for sending and receiving
func assertExampleUDPConnection(ipaddr string) {
	if exampleUDPConnection == nil {
		raddr, errRaddr := net.ResolveUDPAddr("udp", ipaddr+":"+EXAMPLEUDP_PORT)
		logger.WriteError("assertExampleUDPConnection()", errRaddr)

		laddr, errLaddr := net.ResolveUDPAddr("udp", ":"+EXAMPLEUDP_PORT)
		logger.WriteError("assertExampleUDPConnection()", errLaddr)

		var errDial error
		exampleUDPConnection, errDial = net.DialUDP("udp", laddr, raddr)
		logger.WriteError("assertExampleUDPConnection()", errDial)
	}
}

// handle the carried 4 byte data
var ExampleUDPBytesToFloat = make(map[string]func([EXAMPLEUDP_DATA_SIZE]byte) float64)

var ExampleUDPFloatToBytes = make(map[string]func(float64) [EXAMPLEUDP_DATA_SIZE]byte)

func AddExampleUDPBytesToFloat(key string, fn func([EXAMPLEUDP_DATA_SIZE]byte) float64) bool {
	ExampleUDPBytesToFloat[key] = fn
	return true
}

func AddExampleUDPFloatToBytes(key string, fn func(float64) [EXAMPLEUDP_DATA_SIZE]byte) bool {
	ExampleUDPFloatToBytes[key] = fn
	return true
}

// bytes 2 and 3 -> littleEndian uint16 -> float64
func ExampleUDPType1BytesToFloat(b [EXAMPLEUDP_DATA_SIZE]byte) float64 {
	buffer := bytes.NewBuffer(b[2:4])
	var tmp int16
	binary.Read(buffer, binary.LittleEndian, &tmp)

	x := float64(tmp)

	return x
}

var ExampleUDPType1BytesToFloatOk = AddExampleUDPBytesToFloat("Type1", ExampleUDPType1BytesToFloat)

func ExampleUDPType2BytesToFloat(b [EXAMPLEUDP_DATA_SIZE]byte) float64 {
	// Take first two bytes instead of 2 and 3
	buffer := bytes.NewBuffer(b[0:2])
	var tmp uint16
	binary.Read(buffer, binary.LittleEndian, &tmp)

	x := float64(tmp)
	return x
}

var ExampleUDPType2BytesToFloatOk = AddExampleUDPBytesToFloat("Type2", ExampleUDPType2BytesToFloat)

func ExampleUDPType2FloatToBytes(x float64) (b [EXAMPLEUDP_DATA_SIZE]byte) {
	// convert to uint16
	tmp := uint16(x)
	buffer := new(bytes.Buffer)
	err := binary.Write(buffer, binary.LittleEndian, tmp)
	logger.WriteError("ExampleUDPType2FloatToBytes()", err)

	b[0] = buffer.Bytes()[0]
	b[1] = buffer.Bytes()[1]
	fmt.Println(b, x)
	return
}

var ExampleUDPType2FloatToBytesOk = AddExampleUDPFloatToBytes("Type2", ExampleUDPType2FloatToBytes)
