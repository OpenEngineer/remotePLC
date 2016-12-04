package main

import (
	"./graph/"
	"./logger/"
	"./parser/"
	"errors"
	"flag"
	"time"
)

const (
	COMMENT_CHAR       = "#"
	EXTRA_NEWLINE_CHAR = ";"
	LOG_REGEXP         = ".*[^_]$"
)

func main() {
	logger.EventMode = logger.FATAL
	cmdString, fname, timeStep, saveInterval := parseArgs()

	blockTable, lineTable := readTables(cmdString, fname)

	g := graph.ConstructGraph(blockTable, lineTable,
		[]string{"inputs"}, []string{"logic"})

	logger.EventMode = logger.WARNING

	controlLoop(g, timeStep, saveInterval, LOG_REGEXP)
}

func parseArgs() (cmdString, fname string, timeStep time.Duration, saveInterval int) {
	// compile the flags
	cmdStringPtr := flag.String("c", "", "blocks semicolon separated, appended to list of blocks")
	fnamePtr := flag.String("f", "blocks.cfg", "file with list of blocks (default: blocks.cfg)")
	t := flag.String("t", "250ms", "length of cycle in [ms]")
	s := flag.Int("s", 4, "save interval in number of cycles")

	// TODO: add more flags
	flag.Parse()

	// convert to correct datatypes
	cmdString = *cmdStringPtr
	fname = *fnamePtr

	var timeErr error
	timeStep, timeErr = time.ParseDuration(*t)
	logger.WriteError("parseArgs()", timeErr)

	saveInterval = *s

	return
}

func readTables(cmdString, fname string) (groupedBlockTable map[string]([][]string), lineTable [][]string) {

	var t_ parser.ConstructorTable
	t := &t_

	t.ReadAppendFile(fname, []string{"\n", ";"})
	t.ReadAppendString(cmdString, []string{"\n", ";"})

	// add or replace
	t.MergeRows("\\")
	t.WordToLine("|")
	t.SubstituteSingleWordLine("|", [][]int{
		[]int{0, 0},
		[]int{-1, 0}, // name of previous block
		[]int{1, 0},  // name of next block
	}, []string{"Line"})
	t.AddRow([]string{"_", "Node"})
	t.GenerateMissingNames(0, ".*Line$")

	// clean
	t.RemoveComments("#")
	t.RemoveEmptyRows(0)
	t.CorrectSuffixes("_", 2)
	t.DetectDuplicates(0)

	// create the sub tables, and leave the remainder in the block table
	groupedBlockTable["inputs"] = t.FilterTable(1, ".*Input$")
	groupedBlockTable["outputs"] = t.FilterTable(1, ".*Output$")
	groupedBlockTable["logic"] = t.FilterTable(1, ".*Logic$")
	groupedBlockTable["nodes"] = t.FilterTable(1, ".*Node$")
	groupedBlockTable["stops"] = t.FilterTable(1, ".*Stop$")

	lineTable = t.FilterTable(1, ".*Line$")

	// if the constructorTable isnt empty now, then there is a problem
	// TODO: function in parser
	for _, row := range *t {
		logger.WriteError("readTables()",
			errors.New("unknown block type: "+row[1]))
	}

	return
}
