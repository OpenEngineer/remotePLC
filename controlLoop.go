package main

import (
	"./blocks/"
	"./logger/"
  //"fmt"
	"log"
	"os"
  "sort"
	"time"
)

func controlLoop(inputs, outputs, logic, stoppers, lines map[string]blocks.Block) {
	// Preprocessing step, make the orderedLines map
  logger.EventMode = logger.FATAL
	orderedLines := orderLines(lines, inputs, logic) // outputs and terminators are not needed for this

	checkConnectivity(inputs, outputs, logic, lines, orderedLines) // all outputs and all logic must be connected

	// Init the datalogging
	dataLogger := logger.MakeDataLogger()

	// Cycle the input updates in the background
  logger.EventMode = logger.WARNING
	cycleInputs(inputs, 250, 10) // same rate

	// Main loop
	ticker := time.NewTicker(500 * time.Millisecond)
	for {
		<-ticker.C

		for _, v := range orderedLines {
      //fmt.Println("processing line: ", v)
			lines[v].Update()
		}
		for _, v := range logic {
			v.Update()
		}
		for _, v := range orderedLines {
			lines[v].Update()
		}

		cycleOutputs(outputs)

		cycleStoppers(stoppers)

		logData(dataLogger, inputs, outputs, logic)
	}

	logData(dataLogger, inputs, outputs, logic)
}

func orderLines(lines, inputs, logic map[string]blocks.Block) []string {
	// Initialize a counting data structure
	numLines := len(lines)
	count := make(map[string][]int)
	for k, _ := range lines {
		count[k] = make([]int, numLines+1)
	}

	// Initialize the lines to empty arrays
	for _, v := range lines {
		v.Put([]float64{})
	}

	// Initialize the inputs
	for _, v := range inputs {
		v.Put([]float64{1})
	}

	// Initialize the logic and do the logic
	for _, v := range logic {
		v.Put([]float64{1})
		v.Update()
	}

	// Run each of the lines "numLines+1" times
	for i := 0; i < numLines+1; i++ {
		for k, v := range lines {
			v.Update()
			count[k][i] = len(v.Get())
		}
	}

	// List the lines that complete at each level
	complete := make([][]string, numLines)
	for k, _ := range lines {
		for i := 0; i < numLines; i++ {
			if count[k][i] == count[k][i+1] {
				complete[i] = append(complete[i], k)
				break
			}

			if i == numLines {
				log.Fatal("in orderLines(), \"", k, "\", error: circularity detected")
			}
		}
	}

	// Flatten this list to create the final order
	orderedLines := []string{}
	for _, v := range complete {
		orderedLines = append(orderedLines, v...)
	}

	if len(orderedLines) != len(lines) {
		log.Fatal("in orderLines(), error: orderedLined not of same length as lines, \"", orderedLines, "\" vs \"", lines, "\"")
	}
	return orderedLines
}

func checkConnectivity(inputs, outputs, logic, lines map[string]blocks.Block, orderedLines []string) {
	// Initialize the inputs
	for _, v := range inputs {
		v.Put([]float64{1})
	}

	// Initialize the outputs
	for _, v := range outputs {
		v.Put([]float64{})
	}

	// Initialize the logic and do the logic, and then reset
	for _, v := range logic {
		v.Put([]float64{1})
		v.Update()
		v.Put([]float64{})
	}

	// Update the lines
	for _, v := range orderedLines {
		lines[v].Update()
	}

	// Check bad outputs
	for k, v := range outputs {
		if len(v.Get()) == 0 {
			log.Fatal("in checkConnectivity(), output \"", k, "\", error: bad connectivity")
		}
	}

  // not necessarily fatal
	for k, v := range logic {
		v.Update()

		if len(v.Get()) == 0 {
      logger.WriteEvent("in checkConnectivity(), logic \"", k, "\", error: bad connectivity")
		}
	}
}

func cycleInputs(inputs map[string]blocks.Block, period time.Duration, desync time.Duration) {
	for _, v := range inputs {
		time.Sleep(desync * time.Millisecond)
		go func(b blocks.Block) {
			ticker := time.NewTicker(period * time.Millisecond)
			for {
				<-ticker.C
				b.Update()
			}
		}(v)
	}
}

func cycleOutputs(outputs map[string]blocks.Block) {
	ch := make(chan int, len(outputs))

	//fork
	for _, v := range outputs {
		go func(b blocks.Block) {
			b.Update()
			ch <- 0

		}(v)
	}

	// join
	for _, _ = range outputs {
		<-ch
	}
}

func cycleStoppers(stoppers map[string]blocks.Block) {
	isConverged := true
	for k, stopper := range stoppers {
		stopper.Update()
		x := stopper.Get()
		for _, v := range x {
			if v > -1.0 {
				isConverged = false
			}

			if v >= 1.0 {
				log.Fatal("in cycleStoppers, \"", k, "\", error: divergence detected")
			}
		}
	}

	if isConverged {
		os.Exit(0)
	}
}

// TODO: a lexicographical sort operation
func logData(dataLogger *logger.DataLogger, dataSets ...map[string]blocks.Block) {
	fields := []string{}
	data := [][]float64{}

	for _, dataSet := range dataSets { // eg. inputs, outputs
    // first sort the list of keys
    keys := []string{}
    for key, _ := range dataSet {
      keys = append(keys, key)
    }

    sort.Strings(keys)

    // then access these keys to append to the data slice
    for _, key := range keys {
			fields = append(fields, key)
			data = append(data, dataSet[key].Get())
    }
	}

	dataLogger.WriteData(fields, data)
}
