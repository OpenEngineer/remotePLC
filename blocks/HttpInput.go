package blocks

import (
	"../logger/"
	//"bufio"
	"errors"
	"fmt"
	"net/http"
	//"os"
	"strconv"
	"strings"
	"time"
)

type HttpInput struct {
	InputBlockData
	Server       *http.Server
	tmp          []float64
	numInput     int
	numOutput    int
	addTimeStamp bool // if true then downstream floats have one more then upstream floats
	// this timeStamp float is not the actual time, it is just marks that a message has been received (so that cached floats can be distinguished from new floats)
}

func (b *HttpInput) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	urlFields := strings.Split(r.URL.Path, "/")

	if len(urlFields) > 1 {
		lastField := urlFields[len(urlFields)-1]
		numberFields := strings.Split(lastField, ",")

		// distrust the network users and ignore the message if it doesn't match the number of expected fields
		//  DoD attacks are always possible
		if len(numberFields) != b.numInput {
			fmt.Fprintf(w, "error: bad number of inputs")
		} else {
			for i, v := range numberFields {
				number, parseErr := strconv.ParseFloat(v, 64)
				if parseErr == nil {
					b.tmp[i] = number
				} else {
					fmt.Fprintf(w, "error: bad number")
					return
				}
			}

			if b.addTimeStamp {
				b.tmp[b.numOutput-1] = 1.0 - (b.tmp[b.numOutput-1]) // switches between 0 and 1
			}
			numberStr := fmt.Sprintln(b.tmp)
			fmt.Fprintf(w, "%s", numberStr)
		}
	} else {
		fmt.Fprintf(w, "error: bad url")
	}
}

func (b *HttpInput) Update() {
	if len(b.tmp) == b.numInput {
		copy(b.tmp[0:b.numInput], b.in)
		copy(b.tmp, b.out)
	} else {
		logger.WriteEvent("HttpInput, bad number of inputs ", len(b.tmp), " should be ", b.numInput)
		if len(b.out) != b.numOutput {
			b.out = make([]float64, b.numOutput)
		}
	}
}

func HttpInputConstructor(name string, words []string) Block {
	if len(words) != 2 {
		logger.WriteError("HttpInputConstructor()", errors.New("need at least 2 words"))
	}
	port := ":" + words[0]
	numInput, err := strconv.ParseInt(words[1], 10, 64)
	logger.WriteError("HttpInputConstructor()", err)

	addTimeStamp := false
	numOutput := numInput
	if len(words) > 2 {
		if words[2] == "addTimeStamp" {
			addTimeStamp = true
			numOutput = numInput + 1
		} else {
			logger.WriteError("HttpInputConstructor()", errors.New(words[2]+" not a recognized option"))
		}
	}

	b := &HttpInput{
		numInput:     int(numInput),
		tmp:          make([]float64, numOutput), // store incoming data here is it doesn't interfere with the internal Get function
		addTimeStamp: addTimeStamp,
		numOutput:    int(numOutput),
	}
	b.Server = &http.Server{
		Addr:           port,
		Handler:        b,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go b.Server.ListenAndServe()

	return b
}

var HttpInputConstructorOk = AddConstructor("HttpInput", HttpInputConstructor)
