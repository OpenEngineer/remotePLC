package main

import (
	"./blocks/"
	"./logger/"
	//"bufio"
	"errors"
	"flag"
	"os"
	"regexp"
	"strings"
	"time"
)

const (
	COMMENT_CHAR       = "#"
	EXTRA_NEWLINE_CHAR = ";"
)

func main() {
	logger.EventMode = logger.FATAL
	inputTable, outputTable, logicTable, nodeTable, stoppersTable, lineTable, timeStep, saveInterval := readConfig()

	logger.WriteEvent("Constructing inputs:")
	inputs := blocks.ConstructAll(inputTable)
	logger.WriteEvent("Construcing outputs:")
	outputs := blocks.ConstructAll(outputTable)
	logger.WriteEvent("Constructing logic:")
	logic := blocks.ConstructAll(logicTable)
	logger.WriteEvent("Constructing nodes:")
	nodes := blocks.ConstructAll(nodeTable)
	logger.WriteEvent("Constructing stoppers:")
	stoppers := blocks.ConstructAll(stoppersTable)
	logger.WriteEvent("Construction lines:")
	lines := blocks.ConstructAll(lineTable)

	// TODO: add loop time parameters
	logger.EventMode = logger.WARNING
	controlLoop(inputs, outputs, logic, nodes, stoppers, lines, timeStep, saveInterval)
}

func readConfig() (inputTable, outputTable, logicTable, nodeTable, stopTable, lineTable [][]string,
	timeStep time.Duration, saveInterval int) {
	// Compile the flags
	cmdString := flag.String("c", "", "blocks semicolon separated, appended to list of blocks")
	fname := flag.String("f", "blocks.cfg", "file with list of blocks (default: blocks.cfg)")
	t := flag.String("t", "250ms", "length of cycle in [ms]")
	s := flag.Int("s", 4, "save interval in number of cycles")

	// TODO: add more flags
	flag.Parse()

	// now read the files
	blockTable := readFileTable(*fname)
	blockTable = append(blockTable, readStringTable(*cmdString)...)

	// create the sub tables, and leave the remainder in the block table
	inputTable = filterTable(&blockTable, ".*Input$")
	outputTable = filterTable(&blockTable, ".*Output$")
	logicTable = filterTable(&blockTable, ".*Logic$")
	nodeTable = filterTable(&blockTable, ".*Node$")
	lineTable = filterTable(&blockTable, ".*Line$")
	stopTable = filterTable(&blockTable, ".*Stop$")

	// if the blockTable isnt empty now, then there is a problem
	for _, record := range blockTable {
		logger.WriteError("readConfig()",
			errors.New("unknown block type: "+record[1]))
	}

	var timeErr error
	timeStep, timeErr = time.ParseDuration(*t)
	logger.WriteError("readConfig()", timeErr)

	saveInterval = *s

	return
}

func filterTable(tableIn *[][]string, typeRegexp string) (tableOut [][]string) {
	re := regexp.MustCompile(typeRegexp)

	var tmpTable [][]string
	for _, record := range *tableIn {
		if re.MatchString(record[1]) {
			tableOut = append(tableOut, record)
		} else {
			tmpTable = append(tmpTable, record)
		}
	}

	*tableIn = tmpTable

	return
}

func readFileTable(fname string) (table [][]string) {
	file, err := os.Open(fname)
	defer file.Close()
	logger.WriteError("readFileTable()", err)

	// get the size of the file
	finfo, errStat := os.Stat(fname)
	logger.WriteError("readFileTable()", errStat)

	fileBytes := make([]byte, int(finfo.Size()))

	// read the whole file into memory
	_, errRead := file.Read(fileBytes)
	logger.WriteError("readFileTable()", errRead)

	// now process the string with the lower level readStringTable() function
	table = readStringTable(string(fileBytes))

	return
}

func readStringTable(str string) (table [][]string) {
	// first split by newline
	split0 := strings.Split(str, "\n")

	// then loop these and split by semicolon
	var split1 [][]string
	for _, s := range split0 {
		split1 = append(split1, strings.Split(s, EXTRA_NEWLINE_CHAR))
	}

	// Now flatten the first dimension (so that lines split by "\n" and ";" become equivalent)
	var split2 []string
	for _, s := range split1 {
		split2 = append(split2, s...)
	}

	// now split each line-string into its words
	for _, s := range split2 {
		table = append(table, strings.Fields(s))
	}

	table = cleanTable(table)

	return
}

func cleanTable(table [][]string) [][]string {
	table = removeComments(table)
	table = removeEmptyRows(table)

	return table
}

func removeComments(table [][]string) (tableOut [][]string) {
	for _, row := range table {
		rowOut := []string{}
		for _, word := range row {
			if string(word[0]) == COMMENT_CHAR {
				break
			}
			rowOut = append(rowOut, word)
		}

		tableOut = append(tableOut, rowOut)
	}

	return
}

func removeEmptyRows(table [][]string) (tableOut [][]string) {
	for _, row := range table {
		if len(row) > 0 {
			tableOut = append(tableOut, row)
		}
	}

	return
}
