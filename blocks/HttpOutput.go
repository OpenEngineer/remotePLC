package blocks

import (
	"../logger/"
	"fmt"
	"net/http"
)

type HttpOutput struct {
	OutputBlockData
	uriBase string
	prevUri string
}

func (b *HttpOutput) Update() {
	uriSuffix := ""
	for _, v := range b.in {
		uriSuffix = uriSuffix + fmt.Sprintf("%f,", v)
	}

	var uri string
	if string(b.uriBase[len(b.uriBase)-1]) == "/" {
		uri = b.uriBase + uriSuffix
	} else {
		uri = b.uriBase + "/" + uriSuffix
	}

	if uri != b.prevUri {

		_, responseErr := http.Get(uri)
		logger.WriteError("HttpOutput.Update()", responseErr)

		b.out = b.in

		b.prevUri = uri
	}
}

func HttpOutputConstructor(name string, words []string) Block {
	uriBase := words[0]

	// get the general state
	b := &HttpOutput{uriBase: uriBase}
	return b
}

var HttpOutputOk = AddConstructor("HttpOutput", HttpOutputConstructor)
