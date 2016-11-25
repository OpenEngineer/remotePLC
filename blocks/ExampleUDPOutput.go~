package blocks

import (
  "../logger/"
  "errors"
  "fmt"
  "net"
  "log"
  "strconv"
  "strings"
)

// see ExampleUDP.go for protocol details

type ExampleUDPOutput struct {
  OutputBlockData
  nRecords int // fixed, the nRecords value in question is varied
  question ExampleUDPPacket // converted to []byte in Update()

  converter func(float64)[4]byte

  conn net.Conn
}

func (b *ExampleUDPOutput) Update() {
  // If the b.in contains less elements than nRecords, then nRecords is set to len(b.in)
  nRecords := b.nRecords
  if len(b.in) < b.nRecords {
    nRecords = len(b.in)
  }

  // Now set the nSdos in question
  b.question.Header1.Uint1 = uint16(nRecords)

  // loop the input and store the records
  for i, v := range b.in {
    b.question.Records[i].DwData = b.converter(v)
  }

  questionBytes := exampleUDPPacketToBytes(b.question)

  _, errWrite := b.conn.Write(msg)
  logger.WriteError("ExampleUDPOutput.Update()", errWrite)

  answerBytes := make([]byte, 1460)
  _, err := b.conn.Read(answerBytes)
  logger.WriteError("ExampleUDPOutput.Update()", err)

  // now parse the return message
  answer := exampleUDPBytesToPacket(answerBytes)
  // TODO: check the status of the answer

  b.out = b.in[0:nRecords]
  return
}

func ExampleUDPOutputConstructor(words []string) Block {
  ipaddr := words[0]
  converterType := words[1]
  recordWords := words[2:]

  question := ExampleUDPPacket{
    Header1: ExampleUDPHeader1{
      Byte2: 0xbb,
    },
    Header2: ExampleUDPHeader2{
      Uint2: 10,
    },
  }

  // store all the records
  nRecords := len(recordWords)
  if nRecords == 0 {
    log.Fatal("must specify at least one record")
  } else if nRecords > 10 {
    log.Fatal("too many records")
  }

  question.Header2.Uint1 = uint16(nRecords)
  for i, w := range recordWords {
    id, errId := strconv.ParseUint(w, 0, 16)
    logger.WriteError("ExampleUDPOutputConstructor()", errId)

    question.Records[i] = ExampleUDPRecord{
      Uint2: uint16(id),
    }
  }

  assertExampleUDPConnection(ipaddr)

  converter, ok := ExampleUDPFloatToBytes[converterType]
  if !ok {
    logger.WriteError("ExampleUDPOutputConstructor()", errors.New("couldnt find "+converterType)) // TODO: also do this for the Input constructor
  }

  b := &ExampleUDPOutput{
    question: question,
    conn: exampleUDPConnection, 
    converter: converter,
    nRecords: nRecords,
  }

  return b
}

var ExampleUDPOutputConstructorOk = AddConstructor("ExampleUDPOutput", ExampleUDPOutputConstructor)
