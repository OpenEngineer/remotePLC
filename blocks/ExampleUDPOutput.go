package blocks

import (
	"../logger/"
	"errors"
	"net"
	"strconv"
)

// see ExampleUDP.go for protocol details

type ExampleUDPOutput struct {
	OutputBlockData
	numRecords int              // fixed here, because the numRecords value in question is varied
	question   ExampleUDPPacket // filled and converted to []byte in Update()
	dataConv   func(float64) [EXAMPLEUDP_DATA_SIZE]byte
	conn       net.Conn
}

func (b *ExampleUDPOutput) Update() {
	// If the b.in contains less elements than numRecords, then numRecords is set to len(b.in)
	numRecords := b.numRecords
	if len(b.in) < b.numRecords {
		numRecords = len(b.in)
	}

	// Now set the nSdos in question
	b.question.Header.NumRecords = uint16(numRecords)

	// loop the input and store the records
	for i, v := range b.in {
		b.question.Records[i].Data = b.dataConv(v)
	}

	questionBytes := exampleUDPPacketToBytes(b.question)

	_, errWrite := b.conn.Write(questionBytes)
	logger.WriteError("ExampleUDPOutput.Update()", errWrite)

	answerBytes := make([]byte, EXAMPLEUDP_MAX_BYTES)
	_, err := b.conn.Read(answerBytes)
	logger.WriteError("ExampleUDPOutput.Update()", err)

	// now parse the return message
	answer := exampleUDPBytesToPacket(answerBytes) // TODO: create the error herein
	if answer.Header.ErrorCode != 0 {
		logger.WriteError("ExampleUDPOutput.Update()",
			errors.New("protocol error "+strconv.Itoa(int(answer.Header.ErrorCode))))
	}

	b.out = b.in[0:numRecords]
	return
}

func ExampleUDPOutputConstructor(words []string) Block {
	ipaddr := words[0]
	dataConvType := words[1]
	recordIds := words[2:]

	question := ExampleUDPPacket{
		Header: ExampleUDPHeader{
			Version: EXAMPLEUDP_VERSION,
			OpId:    EXAMPLEUDP_WRITEOP,
			SeqNo:   uint16(1),
		},
	}

	// store all the records
	numRecords := len(recordIds)
	if numRecords == 0 {
		logger.WriteError("ExampleUDPOutputConstructor()",
			errors.New("must specify at least one record"))
	} else if numRecords > EXAMPLEUDP_MAX_RECORDS {
		logger.WriteError("ExampleUDPOutputConstructor()",
			errors.New("too many records"))
	}

	question.Header.NumRecords = uint16(numRecords)
	for i, idStr := range recordIds {
		id, errId := strconv.ParseUint(idStr, 0, 16)
		logger.WriteError("ExampleUDPOutputConstructor()", errId)

		question.Records[i] = ExampleUDPRecord{
			Id: uint16(id),
		}
	}

	assertExampleUDPConnection(ipaddr)

	dataConv, ok := ExampleUDPFloatToBytes[dataConvType]
	if !ok {
		logger.WriteError("ExampleUDPOutputConstructor()",
			errors.New("couldnt find "+dataConvType)) // TODO: also do this for the Input constructor
	}

	b := &ExampleUDPOutput{
		question:   question,
		dataConv:   dataConv,
		numRecords: numRecords,
		conn:       exampleUDPConnection,
	}

	return b
}

var ExampleUDPOutputConstructorOk = AddConstructor("ExampleUDPOutput", ExampleUDPOutputConstructor)
