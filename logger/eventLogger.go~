package logger

import "log"
import "os"

type EventModeType int

const (
	QUIET EventModeType = iota
	WARNING
	FATAL
)

var EventMode EventModeType = WARNING

func createEventLogger() *log.Logger {
	fname := "events.log"
	file, err := os.Create(fname)

	if err != nil {
		log.Fatal("in createLogFile(), \"", fname, "\", ", err)
	}

	logger := log.New(file, "", log.Ldate|log.Ltime)
	return logger
}

var events *log.Logger

func WriteEvent(s ...interface{}) {
	if EventMode != QUIET {
		if events == nil {
			events = createEventLogger()
		}

		events.Print(s...)
	}
}

func WriteError(context string, err error) {
	if EventMode != QUIET && err != nil {
		if events == nil {
			events = createEventLogger()
		}

		if EventMode == WARNING {
			events.Print(context, err)
		} else if EventMode == FATAL {
			// quit the program
			log.Fatal(context+", ", err)
			os.Exit(1) // TODO: better error codes?
		} else {
			// TODO: make redundant by using fixed set commands for EventMode
			log.Fatal("event mode not recognized")
		}
	}
}
