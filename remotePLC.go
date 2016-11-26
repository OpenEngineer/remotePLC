package main

import (
	"./blocks/"
	"./logger/"
	"bufio"
	"errors"
	"flag"
	"os"
	"regexp"
	"strings"
)

func main() {
	logger.EventMode = logger.FATAL
	inputTable, outputTable, logicTable, stoppersTable, lineTable := readConfig()

	inputs := blocks.ConstructAll(inputTable)
	outputs := blocks.ConstructAll(outputTable)
	logic := blocks.ConstructAll(logicTable)
	stoppers := blocks.ConstructAll(stoppersTable)
	lines := blocks.ConstructAll(lineTable)

	// TODO: add loop time parameters
	logger.EventMode = logger.WARNING
	controlLoop(inputs, outputs, logic, stoppers, lines)
}

func readConfig() (inputTable, outputTable, logicTable, stopTable, lineTable [][]string) {
	// Compile the flags
	cmdString := flag.String("c", "", "blocks semicolon separated, appended to list of blocks")
	fname := flag.String("f", "blocks.cfg", "file with list of blocks (default: blocks.cfg)")

	// TODO: add more flags
	flag.Parse()

	// now read the files
	blockTable := readFileTable(*fname)
	blockTable = append(blockTable, readStringTable(*cmdString)...)

	// create the sub tables, and leave the remainder in the block table
	inputTable = filterTable(&blockTable, ".*Input$")
	outputTable = filterTable(&blockTable, ".*Output$")
	logicTable = filterTable(&blockTable, ".*(Node|Logic)$")
	lineTable = filterTable(&blockTable, ".*Line$")
	stopTable = filterTable(&blockTable, ".*Stop$")

	// if the blockTable isnt empty now, then there is a problem
	for _, record := range blockTable {
		logger.WriteError("readConfig()",
			errors.New("unknown block type: "+record[1]))
	}

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
	logger.WriteError("readFileTable()", err)

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
			if string(word[0]) == "#" {
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
