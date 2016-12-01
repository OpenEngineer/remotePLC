package lines

import (
	"../blocks/"
)

type Line interface {
	Transfer()
}

type LineData struct {
	b0 []string
	b1 []string
}

var Lines []Line // TODO: or map?

var Constructors = make(map[string]func([]string) Line)

func AddConstructor(key string, fn func([]string) Line) bool {
	Constructors[key] = fn
	return true
}

func Construct(constructorType string, words []string) Line {
	var l Line
	if constructor, ok := Constructors[constructorType]; ok {
		l = constructor(words)
	} else {
		log.Fatal("invalid line constructor: ", constructorType)
	}
	return l
}

func ConstructGlobal(words []string) Line {
	defer logger.WriteEvent("constructed line: ", words)
	l := Construct(words[0], words[1:])
	Lines = append(Lines, l)
	return l
}

func ConstructAll(wordsTable [][]string) []Line {
	lines := []Line{}

	for _, words := range wordsTable {
		lines = append(lines, ConstructGlobal(words))
	}

	return lines
}

func PrepareLines(inputs, outputs, logic map[string]blocks.Block) {
	blocks.BlockMode = blocks.CONNECTIVITY
	orderLines(inputs, logic)

	checkConnectivity(inputs, outputs, logic)
	blocks.BlockMode = blocks.REGULAR
}

func orderLines(inputs, logic map[string]blocks.Block) {

	// Initialize a counting data structure
	numLines := len(Lines)
	count := make(map[int][]int)
	for i, _ := range Lines {
		count[i] = make([]int, numLines+1)
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
		for j, v := range Lines {
			v.Transfer()
			count[j][i] = len(v.Get()) // TODO: figure this out
		}
	}

	// List the lines that complete at each level
	complete := make([][]int, numLines)
	for i, _ := range Lines {
		for j := 0; j < numLines; j++ {
			if count[i][j] == count[i][j+1] {
				complete[j] = append(complete[j], i)
				break
			}

			if j == numLines {
				log.Fatal("in orderLines(), \"", i, "\", error: circularity detected")
			}
		}
	}

	// Flatten this list to create the final order
	orderedLines := []Line{}
	for _, v := range complete {
		for _, w := range v {
			orderedLines = append(orderedLines, Lines[v])
		}
	}

	if len(orderedLines) != len(Lines) {
		log.Fatal("in orderLines(), error: orderedLined not of same length as lines, \"", orderedLines, "\" vs \"", Lines, "\"")
	}

	Lines = orderedLines[:]
}

func checkConnectivity(inputs, outputs, logic map[string]blocks.Block) {
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
	for _, l := range Lines {
		l.Transfer() // should fix the nodes
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
