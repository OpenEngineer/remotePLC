package blocks

import (
	"../logger/"
	"bufio"
	"errors"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

// implemented as copy of TimeFileInput
// TODO: factor out common functionality

type MapNode struct {
	BlockData
	fname   string
	file    *os.File // usefull for seeking
	modTime time.Time
	xInterp []float64
	yInterp [][]float64

	// options
	cycle    bool
	discrete bool // do not interpolate (nearest behaviour)
}

// TODO: log error messages
func checkMonotonicity_(xInterp []float64) bool {
	isMonotone := true

	for i, x := range xInterp[1:] {
		if x < xInterp[i] {
			isMonotone = false
		}
	}

	return isMonotone
}

// Example of runtime file modification
func (b *MapNode) readFile() error {
	if b.file == nil {
		return errors.New("error: MapNode nil")
	}
	b.file.Seek(0, 0) // also seeks to 0 when piping stdin

	scanner := bufio.NewScanner(b.file) // use default Split: ScanLines

	logger.WriteEvent("reading MapNode file ", b.fname)

	xInterpNew := []float64{}
	yInterpNew := [][]float64{}

	for scanner.Scan() {
		line := scanner.Text()
		words := strings.Fields(line)

		// parse the time
		xInterp, _ := strconv.ParseFloat(words[0], 64)
		xInterpNew = append(xInterpNew, xInterp)

		// parse the values
		yInterp := []float64{}
		for _, word := range words[1:] {
			y, _ := strconv.ParseFloat(word, 64)
			yInterp = append(yInterp, y)
		}
		yInterpNew = append(yInterpNew, yInterp)
	}

	// change the block only if the values are monotone
	if len(xInterpNew) == 0 {
		return errors.New("file NOK")
	} else if checkMonotonicity_(xInterpNew) {
		b.xInterp = xInterpNew
		b.yInterp = yInterpNew
		info, _ := b.file.Stat()
		b.modTime = info.ModTime()
	} else {
		return errors.New("warning: invalid MapNode, ignoring")
	}

	return nil
}

func (b *MapNode) refreshFile() {
	// reread the time file
	var openErr error
	b.file, openErr = reopenFile_(b.fname, b.file)
	if openErr != nil {
		logger.WriteEvent("warning ignoring timefile ", b.fname, ", ", openErr)
		return
	}

	info, statErr := b.file.Stat()
	if statErr != nil {
		logger.WriteEvent("warning ignoring timefile ", b.fname, ", ", statErr)
		return
	}

	if !b.modTime.Equal(info.ModTime()) {
		readErr := b.readFile()
		if readErr != nil {
			logger.WriteEvent("warning ignoring timefile ", readErr)
		}
	}
}

func (b *MapNode) findInterpSlice(x float64) (int, int, float64) {
	// find the lower index
	if b.cycle {
		xPeriod := b.xInterp[len(b.xInterp)-1] - b.xInterp[0]
		if x > xPeriod {
			x = x - xPeriod*float64(int((x-b.xInterp[0])/xPeriod))
		} else if x < b.xInterp[0] {
			x = x + xPeriod*float64(int((x-b.xInterp[0])/xPeriod)+1)
		}

		if x > xPeriod { // TODO: check and remove this
			log.Fatal("bad calc in findInterpSlice()")
		}
	}

	var iLower int
	var iUpper int
	if x < b.xInterp[0] {
		iLower = 0
		iUpper = 0
	} else {
		for iLower = 0; iLower < len(b.xInterp)-1; iLower++ {
			if b.xInterp[iLower] <= x && x < b.xInterp[iLower+1] {
				break
			}
		}

		// upper index
		iUpper = iLower + 1
		if iLower == len(b.xInterp)-1 {
			iUpper = iLower
		}
	}

	// interpolation factor
	dx := b.xInterp[iUpper] - b.xInterp[iLower]
	var alpha float64
	if dx != 0.0 {
		alpha = (x - b.xInterp[iLower]) / dx
	} else {
		alpha = 1.0
	}

	return iLower, iUpper, alpha
}

func (b *MapNode) Put(x []float64) {
	if len(x) == 0 {
		return // don't do anything
	}
	b.refreshFile()

	// find the time index
	iLower, iUpper, alpha := b.findInterpSlice(x[0])
	y0 := b.yInterp[iLower]
	y1 := b.yInterp[iUpper]

	var yInterp []float64

	// interpolate
	if !b.discrete {
		if len(y0) < len(y1) {
			y0 = append(y0, make([]float64, len(y1)-len(y0))...)
		} else if len(y1) < len(y0) {
			y1 = append(y1, make([]float64, len(y0)-len(y1))...)
		}

		yInterp = make([]float64, len(y0))

		for i := 0; i < len(yInterp); i++ { // in place modification
			yInterp[i] = (1.0-alpha)*y0[i] + alpha*y1[i]
		}
	} else if alpha < 0.5 {
		yInterp = y0
	} else {
		yInterp = y1
	}

	b.in = yInterp
	b.out = yInterp
}

func reopenFile_(fname string, file *os.File) (*os.File, error) {
	var err error

	if file != nil {
		file.Close()
	}

	if fname == "stdin" {
		file = os.Stdin
		err = nil
	} else {
		file, err = os.Open(fname)
	}

	return file, err
}

func MapNodeConstructor(name string, words []string) Block {
	var file *os.File
	fname := words[0]
	var err error
	file, err = reopenFile_(fname, file)
	if err != nil {
		log.Fatal(err)
	}

	// todo: much better argument parser
	cycle := false
	discrete := false
	if len(words) > 1 {
		if words[1] == "Cycle" {
			cycle = true
		} else if words[1] == "Discrete" {
			discrete = true
		}
	}

	if len(words) > 2 {
		if words[1] == "Cycle" || words[2] == "Cycle" {
			cycle = true
		}

		if words[1] == "Discrete" || words[2] == "Discrete" {
			discrete = true
		}
	}

	if cycle == true {
		logger.WriteEvent("cycling MapNode")
	}

	b := &MapNode{fname: fname, file: file, cycle: cycle, discrete: discrete}

	readErr := b.readFile()
	logger.WriteError("error in MapNodeConstructor: "+fname+" not in valid format", readErr)

	logger.WriteEvent("constructed MapNode " + name)
	return b
}

var MapNodeConstructorOk = AddConstructor("MapNode", MapNodeConstructor)
