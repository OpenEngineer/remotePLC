package main

import (
	"./blocks/"
	"./logger/"
	"log"
	"os"
	//"regexp"
	//"sort"
	"time"
)

func controlLoop(inputs, outputs, logic, nodes, stoppers, lines map[string]blocks.Block, timeStep time.Duration,
	saveInterval int) {
	// Preprocessing step, make the orderedLines map
	logger.EventMode = logger.FATAL
	blocks.BlockMode = blocks.CONNECTIVITY
	orderedLines := orderLines(lines, inputs, logic) // outputs and terminators are not needed for this

	checkConnectivity(inputs, outputs, logic, lines, orderedLines) // all outputs and all logic must be connected
	blocks.BlockMode = blocks.REGULAR

	// Precycle the inputs at least once
	for _, v := range inputs {
		v.Update()
	}

	// Cycle the input updates in the background
	logger.EventMode = logger.WARNING

	cycleInputs(inputs, timeStep, 10) // same rate

	// Main loop
	ticker := time.NewTicker(timeStep)
	counter := 0
	for {
		<-ticker.C

		for _, v := range orderedLines {
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

		if counter%saveInterval == 0 {
			blocks.LogData()
			//logData(dataLogger)
			counter = 0
		}
		counter += 1
	}

	// also save the last time
	blocks.LogData()
	//logData(dataLogger)
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
		lines[v].Update() // should fix the nodes
	}

	// Check bad outputs
	for k, v := range outputs {
		if len(v.Get()) == 0 {
			log.Fatal("in checkConnectivity(), output \"", k, "\", error: bad connectivity")
		}
	}

	// not necessarily fatal: TODO: distinction between logic and nodes
	for k, v := range logic {
		v.Update()

		if len(v.Get()) == 0 {
			logger.WriteEvent("in checkConnectivity(), logic \"", k, "\", error: bad connectivity")
			logger.WriteEvent(k, v.Get())
		}
	}
}

func cycleInputs(inputs map[string]blocks.Block, period time.Duration, desync time.Duration) {
	for _, v := range inputs {
		time.Sleep(desync * time.Millisecond)
		go func(b blocks.Block) {
			ticker := time.NewTicker(period)
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
