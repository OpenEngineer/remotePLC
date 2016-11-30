package blocks

import (
	"../logger/"
	"errors"
	"net"
	"strconv"
)

// see ExampleUDP.go for protocol details (data structure, helper functions etc)

type ExampleUDPInput struct {
	InputBlockData
	questionBytes []byte
	dataConv      func([EXAMPLEUDP_DATA_SIZE]byte) float64 // for the carried data
	conn          net.Conn                                 // local copy of exampleUDPConnection
}

func (b *ExampleUDPInput) Update() {
	// send the question message
	_, errWrite := b.conn.Write(b.questionBytes)
	logger.WriteError("ExampleUDPInput.Update()", errWrite)

	// receive the answer message
	answerBytes := make([]byte, EXAMPLEUDP_MAX_BYTES) // max length of a udp packet
	_, errRead := b.conn.Read(answerBytes)
	logger.WriteError("ExampleUDPInput.Update()", errRead)

	// now parse the answer message
	answer := exampleUDPBytesToPacket(answerBytes)

	if answer.Header.ErrorCode != 0 {
		logger.WriteError("ExampleUDPInput.Update()",
			errors.New("protocol error "+strconv.Itoa(int(answer.Header.ErrorCode))))
	} else if errWrite == nil && errRead == nil {
		// Get the number of records and populate b.out
		numRecords := int(answer.Header.NumRecords)
		if len(b.out) != numRecords {
			b.out = make([]float64, numRecords)
		}

		for i := 0; i < numRecords; i++ {
			b.out[i] = b.dataConv(answer.Records[i].Data)
		}

		b.in = b.out
	}
	return
}

func ExampleUDPInputConstructor(name string, words []string) Block {
	ipaddr := words[0]
	dataConvType := words[1]
	recordIds := words[2:]

	question := ExampleUDPPacket{
		Header: ExampleUDPHeader{
			Version: EXAMPLEUDP_VERSION,
			OpId:    EXAMPLEUDP_READOP,
			SeqNo:   uint16(1), // TODO: management of SeqNo
		},
	}

	// get the number of records, and check validity
	numRecords := len(recordIds)
	if numRecords == 0 {
		logger.WriteError("ExampleUDPInputConstructor()",
			errors.New("need at least one record"))
	} else if numRecords > EXAMPLEUDP_MAX_RECORDS {
		logger.WriteError("ExampleUDPInputConstructor()",
			errors.New("too many records specified"))
	}

	// store all the records
	question.Header.NumRecords = uint16(numRecords)
	for i, idStr := range recordIds {
		id, errId := strconv.ParseUint(idStr, 0, 16)
		logger.WriteError("ExampleUDPInputConstructor()", errId)

		question.Records[i] = ExampleUDPRecord{
			Id: uint16(id),
		}
	}

	// convert the message into a []byte
	questionBytes := exampleUDPPacketToBytes(question)

	assertExampleUDPConnection(ipaddr)

	b := &ExampleUDPInput{
		questionBytes: questionBytes,
		conn:          exampleUDPConnection, // local copy
		dataConv:      ExampleUDPBytesToFloat[dataConvType],
	}

	return b
}

var ExampleUDPInputConstructorOk = AddConstructor("ExampleUDPInput", ExampleUDPInputConstructor)
