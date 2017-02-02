package blocks

import (
	//"../logger/"
	"../parser/"
	//"bufio"
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
			b.Update() // assert that b.out has the correct size
			for i, v := range numberFields {
				number, parseErr := strconv.ParseFloat(v, 64)
				if parseErr == nil {
					b.out[i] = number
				} else {
					fmt.Fprintf(w, "error: bad number")
					return
				}
			}

			// send a reply to the client
			numberStr := fmt.Sprintln(b.out)
			fmt.Fprintf(w, "%s", numberStr)
		}
	} else {
		fmt.Fprintf(w, "error: bad url")
	}
}

func (b *HttpInput) Update() {
	if len(b.out) != b.numInput {
		b.out = make([]float64, b.numInput)
	}

	b.in = b.out
}

func HttpInputConstructor(name string, words []string) Block {
	var port string
	var numInput int

	positional := parser.PositionalArgs(&port, &numInput)

	parser.ParsePositionalArgs(words, positional)

	b := &HttpInput{
		numInput: int(numInput),
	}
	b.Server = &http.Server{
		Addr:           ":" + port,
		Handler:        b,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go b.Server.ListenAndServe()

	return b
}

var HttpInputConstructorOk = AddConstructor("HttpInput", HttpInputConstructor)
