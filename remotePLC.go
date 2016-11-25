package main

import (
	"./blocks/"
	"./logger/"
	"bufio"
	"flag"
	"log"
	"os"
	"strings"
)

func main() {
  logger.EventMode = logger.FATAL
	inputTable, outputTable, logicTable, stoppersTable, lineTable := readInput()

	inputs := blocks.ConstructAll(inputTable)
	outputs := blocks.ConstructAll(outputTable)
	logic := blocks.ConstructAll(logicTable)
	stoppers := blocks.ConstructAll(stoppersTable)
	lines := blocks.ConstructAll(lineTable)

	// TODO: add loop time parameters
  logger.EventMode = logger.WARNING
	controlLoop(inputs, outputs, logic, stoppers, lines)
}

func readInput() (inputTable, outputTable, logicTable, stopTable, lineTable [][]string) {
	// Compile the flags
	inputCmd := flag.String("i", "", "inputs semicolon separated")
	inputFname := flag.String("I", "inputs.cfg", "inputs file")
	outputCmd := flag.String("o", "", "outputs semicolon separated")
	outputFname := flag.String("O", "outputs.cfg", "outputs file")
	logicCmd := flag.String("g", "", "logic semicolon separated")
	logicFname := flag.String("G", "logic.cfg", "logic file")
	stopCmd := flag.String("s", "", "stops semicolon separated")
	stopFname := flag.String("S", "stops.cfg", "stops file")
	linesCmd := flag.String("l", "", "logic semicolon separated")
	linesFname := flag.String("L", "lines.cfg", "stops file")

	// TODO: add more flags
	flag.Parse()

	// now read all the files
	inputTable = append(inputTable, readFileTable(*inputFname, "inputs.cfg")...)
	inputTable = append(inputTable, readStringTable(*inputCmd)...)
	outputTable = append(outputTable, readFileTable(*outputFname, "outputs.cfg")...)
	outputTable = append(outputTable, readStringTable(*outputCmd)...)
	logicTable = append(logicTable, readFileTable(*logicFname, "logic.cfg")...)
	logicTable = append(logicTable, readStringTable(*logicCmd)...)
	stopTable = append(stopTable, readFileTable(*stopFname, "stops.cfg")...)
	stopTable = append(stopTable, readStringTable(*stopCmd)...)
	lineTable = append(lineTable, readFileTable(*linesFname, "lines.cfg")...)
	lineTable = append(lineTable, readStringTable(*linesCmd)...)

	return
}

// TODO: ignore comments and empty lines
func readFileTable(fname string, fnameFallback string) (table [][]string) {
	file, err := os.Open(fname)
	if err != nil {
		if fname != fnameFallback {
			log.Fatal(err)
		} else {
			logger.WriteEvent("couldnt find default ", fnameFallback)
			return
		}
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		words := strings.Fields(line)

		table = append(table, words)
	}
	file.Close()

	table = cleanTable(table)

	return
}

func readStringTable(str string) (table [][]string) {
	// first split by newline
	split0 := strings.Split(str, "\n")

	// then loop these and split by semicolon
	var split1 [][]string
	for _, s := range split0 {
		split1 = append(split1, strings.Split(s, ";"))
	}

	// Now flatten the first dimension
	var split2 []string
	for _, s := range split1 {
		split2 = append(split2, s...)
	}

	// now split each string into its tokens
	for _, s := range split2 {
		table = append(table, strings.Fields(s))
	}

	table = cleanTable(table)

	return
}

func cleanTable(table [][]string) [][]string {
	table = deleteComments(table)
	table = deleteEmptyRows(table)

	return table
}

func deleteComments(table [][]string) (tableOut [][]string) {
	for _, row := range table {
		rowOut := []string{}
		for _, word := range row {
			if string(word[0]) == "#" {
				break
			}
			rowOut = append(rowOut, word)
		}

		tableOut = append(tableOut, rowOut)
	}

	return
}

func deleteEmptyRows(table [][]string) (tableOut [][]string) {
	for _, row := range table {
		if len(row) > 0 {
			tableOut = append(tableOut, row)
		}
	}

	return
}
