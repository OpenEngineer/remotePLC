package logger

import "log"
import "fmt"
import "os"
import "strconv"
import "time"
import "regexp"
import "bufio"
import "strings"
import "sort"

var dataHeader string

var dataHeaderPrefix string = "#date time"

// TODO: external input
var maxDataFileSize int64 = 1024

// TODO: external input
var maxTotalDataFilesSize int64 = 1024 * 2

func dataFileRegexp() *regexp.Regexp {
	str := "[0-9]{4}-[0-9]{2}-[0-9]{2}_[0-9]{2}h[0-9]{2}m[0-9]{2}.log"
	re := regexp.MustCompile(str)

	return re
}

//func readdirDataFiles(d string)
func listDataFiles(dirStr string) ([]string, []int64, int64) {
	re := dataFileRegexp()
	dir, _ := os.Open(dirStr)
	infos, _ := dir.Readdir(-1)
	dir.Close()
	totalSize := int64(0)
	fnames := []string{}
	sizes := []int64{}
	for _, info := range infos {
		if re.MatchString(info.Name()) {
			sizes = append(sizes, info.Size())
			totalSize += info.Size()
			fnames = append(fnames, info.Name())
		}
	}
	sort.Strings(fnames)

	return fnames, sizes, totalSize
}

// delete files matching format below in case total size in folder is too large
func cleanDataDir(dirStr string) {
	fnames, sizes, totalSize := listDataFiles(dirStr)

	for i, fname := range fnames {
		if totalSize > maxTotalDataFilesSize {
			// get the first line, remove the prefix
			f, _ := os.Open(fname)
			r := bufio.NewReader(f)
			byteLine0, _, _ := r.ReadLine()
			line0 := string(byteLine0)
			tokens0 := strings.Split(line0, " ")
			if len(tokens0) > 2 {
				tokens0 = tokens0[2:]
			}
			line0 = strings.Join(tokens0, " ")

			// do the same with the last line of data
			f.Seek(0, 2)
			byteLine1, _, _ := r.ReadLine()
			line1 := string(byteLine1)
			tokens1 := strings.Split(line1, " ")
			if len(tokens1) > 2 {
				tokens1 = tokens1[2:]
			}
			line1 = strings.Join(tokens1, " ")

			f.Close()

			WriteEvent("directory overflow, deleting \"", fname, "\":")
			WriteEvent(line0)
			WriteEvent(line1)
			os.Remove(fname)
			totalSize = totalSize - sizes[i]
		} else {
			break
		}
	}

}

// TODO: refactor
func createDataLogger() *log.Logger {
	if dataFile != nil {
		err := dataFile.Close()
		if err != nil {
			log.Fatal("in createDataLogger()", err)
		}
	}

	cleanDataDir(".")

	fname := time.Now().Format("2006-01-02_15h04m05.log")
	file, err := os.Create(fname)
	if err != nil {
		log.Fatal("in createDataLogger(), \"", fname, "\", ", err)
	}

	logger := log.New(file, "", 5)

	dataFile = file

	return logger
}

// TODO: put Data specific package variables into a struct
var dataFile *os.File

var dataLogger *log.Logger = createDataLogger()

func compileDataRecord(fields []string, data [][]float64) (header string, dataStr string) {
	for i, field := range fields {
		d := data[i]

		for i := 0; i < len(d); i++ {
			fieldSuffix := ""
			if len(d) > 1 {
				fieldSuffix = "_" + strconv.FormatInt(int64(i), 10)
			}
			header = header + " " + field + fieldSuffix
			dataStr = dataStr + strconv.FormatFloat(d[i], 'E', -1, 64) + " "
		}
	}

	return header, dataStr
}

func getDataLogger() *log.Logger {
	// is dataFile too big and does a new one need to be created?
	info, _ := dataFile.Stat()
	if info.Size() > maxDataFileSize {
		dataHeader = "" // reset the dataHeader, so that the new proper is determined
		dataLogger = createDataLogger()
	}

	return dataLogger
}

func WriteData(fields []string, x [][]float64) {
	// generate the header and the data
	newDataHeader, dataString := compileDataRecord(fields, x)

	l := getDataLogger()

	// compare the header, if it is different print it
	if newDataHeader != dataHeader {
		dataHeader = newDataHeader

		fmt.Fprintln(dataFile, dataHeaderPrefix+dataHeader)
	}

	// Now print the data
	l.Print(dataString)
}
