package main

import "./blocks/"
import "time"
import "log"
import "os"
import "./logger/"

//import "fmt"

func main() {
	inputTable := [][]string{
		[]string{"var1", "ConstantInput", "666.666"},
		[]string{"var2", "ConstantInput", "0.5"},
		[]string{"var3", "ZeroInput"},
		[]string{"var4", "ScaledInput", "2.0", "1.0", "ConstantInput", "2002"},
		[]string{"var5", "TimeFileInput", "t1.dat"},
		[]string{"var6", "TimeFileInput", "t0.dat"},
	}
	inputs := blocks.ConstructAll(inputTable)

	outputTable := [][]string{
		[]string{"out1", "FileOutput", "out1.dat"},
		[]string{"out2", "PhilipsHueBridgeOutput", "192.168.1.6", "7mRkE0WShaySGxnfL2jQDdMYMXvBAsgtf1n847iA", "2"},
		[]string{"out3", "PhilipsHueBridgeOutput", "192.168.1.6", "7mRkE0WShaySGxnfL2jQDdMYMXvBAsgtf1n847iA", "3"},
		[]string{"out4", "PhilipsHueBridgeOutput", "192.168.1.6", "7mRkE0WShaySGxnfL2jQDdMYMXvBAsgtf1n847iA", "1"},
		//[]string{"out3", "PhilipsHueOutput", "192.168.1.100", "gawdlP-23CzKxbGc6IkNJdwNNSCTCI40y2RbBc-G", "00:17:88:01:10:36:fb:c0-0b"},
	}
	outputs := blocks.ConstructAll(outputTable)

	logicTable := [][]string{
		[]string{"delay1", "DelayLogic"},
	}
	logic := blocks.ConstructAll(logicTable)

	stoppersTable := [][]string{
		[]string{"time", "TimeStop", "3m"},
	}
	stoppers := blocks.ConstructAll(stoppersTable)

	nodeTable := [][]string{
		[]string{"node1", "Node"},
	}
	blocks.ConstructAll(nodeTable)

	lineTable := [][]string{
		//[]string{"line1", "RegexpForkLine", "var3", "out[0-9]"},
		[]string{"line1", "ForkLine", "var5", "out2", "out3"},
		[]string{"line2", "Line", "node1", "out1"},
		[]string{"line3", "Line", "var1", "node1"},
		[]string{"line4", "Line", "var3", "delay1"},
		[]string{"line5", "Line", "var6", "out4"},
	}
	lines := blocks.ConstructAll(lineTable)
	orderedLines := orderLines(lines, inputs, logic)               // outputs and terminators are not needed for this
	checkConnectivity(inputs, outputs, logic, lines, orderedLines) // all outputs and all logic must be connected

	dataLogger := logger.MakeDataLogger()

	cycleInputs(inputs, 50, 1)

	ticker := time.NewTicker(250 * time.Millisecond)
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

		logData(dataLogger, inputs, outputs, logic)

		cycleStoppers(stoppers)
	}
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

	for k, v := range logic {
		v.Update()

		if len(v.Get()) == 0 {
			log.Fatal("in checkConnectivity(), logic \"", k, "\", error: bad connectivity")
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

func logData(dataLogger *logger.DataLogger, dataSets ...map[string]blocks.Block) {
	fields := []string{}
	data := [][]float64{}

	for _, dataSet := range dataSets {
		for k, v := range dataSet {
			fields = append(fields, k)
			data = append(data, v.Get())
		}
	}

	dataLogger.WriteData(fields, data)
}
