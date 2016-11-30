package logger

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"time"
)

// TODO: only one datalogger makes sense
// TODO: remove the header
var DataLoggerInitialized = false

var DataLogger *DataLoggerType

type DataLoggerType struct {
	headerPrefix string
	header       string // including the headerPrefix

	summaryHeader string // including the headerPrefix
	summaryFile   *os.File

	fnameRegexp *regexp.Regexp
	fnameFormat string

	maxFileSize  int64
	maxTotalSize int64

	dir        *os.File
	file       *os.File
	dataLogger *log.Logger
}

// TODO: more arguments
func makeDataLogger() {
	dir, dirErr := os.Open(".")
	if dirErr != nil {
		log.Fatal(dirErr)
	}

	summaryFile, summaryErr := os.Create("summary.log")
	if summaryErr != nil {
		log.Fatal(summaryErr)
	}

	DataLogger = &DataLoggerType{
		headerPrefix: "#date time",
		fnameRegexp:  regexp.MustCompile("^[0-9]{4}-[0-9]{2}-[0-9]{2}_[0-9]{2}h[0-9]{2}m[0-9]{2}.log$"),
		fnameFormat:  "2006-01-02_15h04m05.log",
		maxFileSize:  1024 * 1024 * 10, // 10MB
		maxTotalSize: 1024 * 1024 * 50, // 50MB
		dir:          dir,
		summaryFile:  summaryFile,
	}

	DataLogger.reopenLogger()
}

func (d *DataLoggerType) reopenLogger() {
	if d.file != nil {
		err := d.file.Close()
		if err != nil {
			log.Fatal("in reopenLogger()", err)
		}
	}

	d.cleanDir()

	fname := time.Now().Format(d.fnameFormat)
	file, err := os.Create(fname)
	if err != nil {
		log.Fatal("in createDataLogger(), \"", fname, "\", ", err)
	}

	d.dataLogger = log.New(file, "", 5) // the old one is simply collected as garbage

	d.header = "" // reset the header, so that the new proper one is determined and written

	d.file = file
}

func (d *DataLoggerType) listFiles() (fnames []string, sizes []int64, totalSize int64) {
	d.dir.Seek(0, 0)
	infos, err := d.dir.Readdir(-1)
	if err != nil {
		WriteEvent("warning, problem reading datalog directory: ", err)
	}

	for _, info := range infos {
		if d.fnameRegexp.MatchString(info.Name()) {
			sizes = append(sizes, info.Size())
			totalSize += info.Size()
			fnames = append(fnames, info.Name())
		}
	}
	sort.Strings(fnames)

	return
}

func readCurrentLine(reader *bufio.Reader) (line string, err error) {
	var b []byte
	var isPrefix bool
	b, isPrefix, err = reader.ReadLine()

	if err == nil && isPrefix { // prefer err over "isPrefix"
		err = errors.New("error: couldn't get whole line in readCurrentLine()")
	}

	line = string(b)

	return
}

// variable number of last lines
func readFirstAndLastLine(file *os.File) (line0 string, line1 string, err error) {
	reader := bufio.NewReader(file)

	var err0 error
	line0, err0 = readCurrentLine(reader)

	// do the same with the last line of data
	file.Seek(0, 2)
	var err1 error
	line1, err1 = readCurrentLine(reader)

	// prefer err0 over err1
	if err0 != nil {
		err = err0
	} else if err1 != nil {
		err = err1
	} else {
		err = nil
	}

	return
}

func (d *DataLoggerType) writeSummary(line string) error {
	_, err := fmt.Fprintln(d.summaryFile, line)
	return err
}

// delete files matching format below in case total size in folder is too large
func (d *DataLoggerType) cleanDir() {
	fnames, sizes, totalSize := d.listFiles()

	for i, fname := range fnames {
		if sizes[i] == 0 {
			WriteEvent("removing empty datalog file ", fname)
			os.Remove(fname)
		} else if totalSize > d.maxTotalSize {
			file, err := os.Open(fname)
			if err == nil {
				line0, line1, readErr := readFirstAndLastLine(file)

				if readErr != nil {
					WriteEvent("warning, problem summarizing ", fname, ": ", readErr)
				} else {
					// write the lines to the summary
					if d.summaryHeader != line0 {
						summaryErr := d.writeSummary(line0)
						if summaryErr != nil {
							WriteEvent("warning, with summary file ", summaryErr)
						} else {
							d.summaryHeader = line0
						}
					}

					summaryErr := d.writeSummary(line1)
					if summaryErr != nil {
						WriteEvent("warning, with summary file ", summaryErr)
					}
				}

				file.Close()
			} else {
				WriteEvent("warning, problem summarizing ", fname, ": ", err)
			}

			WriteEvent("summarized \"", fname, "\":")

			os.Remove(fname)
			totalSize = totalSize - sizes[i]
		} else {
			break
		}
	}
}

// TODO: custom record format
func (d *DataLoggerType) compileRecord(fields []string, data [][]float64) (header string, record string) {
	for i := 0; i < len(fields); i++ {
		n := len(data[i])
		for j := 0; j < n; j++ {
			fieldSuffix := ""
			if n > 1 {
				fieldSuffix = "_" + strconv.FormatInt(int64(j), 10)
			}

			// append to header
			header = header + " " + fields[i] + fieldSuffix

			// append to record
			record = record + strconv.FormatFloat(data[i][j], 'E', -1, 64) + " "
		}
	}

	// prepend prefix to header
	header = d.headerPrefix + header

	// the date/time-stamp prefix of the record is added later by the d.dataLogger object

	return header, record
}

func (d *DataLoggerType) refreshLogger(header string) error {
	// is dataLog too big and does a new one need to be created?
	info, err := d.file.Stat()
	if err != nil {
		return err
	}
	if info.Size() > d.maxFileSize || header != d.header {
		d.reopenLogger()
	}

	return nil
}

// the only global function in the dataLogger part of this module
func WriteData(fields []string, x [][]float64) {
	if !DataLoggerInitialized {
		makeDataLogger()
		DataLoggerInitialized = true
	}

	// TODO: all in globals
	d := DataLogger

	// generate the header and the data
	header, record := d.compileRecord(fields, x)

	d.refreshLogger(header)

	// compare the header, if it is different print it
	if header != d.header {
		d.header = header

		fmt.Fprintln(d.file, header)
	}

	// now print the data, prefixed by a date/time-stamp
	d.dataLogger.Print(record)
}
