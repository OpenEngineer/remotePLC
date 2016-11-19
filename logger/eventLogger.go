package logger

import "log"
import "os"

func createEventLogger() *log.Logger {
	fname := "events.log"
	file, err := os.Create(fname)

	if err != nil {
		log.Fatal("in createLogFile(), \"", fname, "\", ", err)
	}

	logger := log.New(file, "", log.Ldate|log.Ltime)
	return logger
}

var events *log.Logger = createEventLogger()

func WriteEvent(s ...interface{}) {
	events.Print(s...)
}
