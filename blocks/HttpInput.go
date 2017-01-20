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
	Server   *http.Server
	tmp      []float64
	numInput int
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

			numberStr := fmt.Sprintln(b.tmp)
			fmt.Fprintf(w, "%s", numberStr)
		}
	} else {
		fmt.Fprintf(w, "error: bad url")
	}
}

func (b *HttpInput) Update() {
	if len(b.tmp) == b.numInput {
		b.in = b.tmp
		b.out = b.tmp
	} else {
		logger.WriteEvent("HttpInput, bad number of inputs ", len(b.tmp), " should be ", b.numInput)
		if len(b.out) != b.numInput {
			b.out = make([]float64, b.numInput)
		}
	}
}

func HttpInputConstructor(name string, words []string) Block {
	if len(words) != 2 {
		logger.WriteError("HttpInputConstructor()", errors.New("need 2 words"))
	}
	port := ":" + words[0]
	numInput, err := strconv.ParseInt(words[1], 10, 64)
	logger.WriteError("HttpInputConstructor()", err)

	b := &HttpInput{
		numInput: int(numInput),
		tmp:      make([]float64, numInput), // store incoming data here is it doesn't interfere with the internal Get function
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
