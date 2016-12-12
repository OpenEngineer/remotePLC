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

type TimeFileInput struct {
	InputBlockData
	fname   string
	file    *os.File // usefull for seeking
	modTime time.Time
	start   time.Time
	tInterp []float64
	xInterp [][]float64

	// options
	cycle    bool
	discrete bool // do not interpolate (nearest behaviour)
}

// TODO: log error messages
func checkMonotonicity(tInterp []float64) bool {
	isMonotone := true

	if tInterp[0] != 0.0 {
		isMonotone = false
	}

	for i, t := range tInterp[1:] {
		if t < tInterp[i] {
			isMonotone = false
		}
	}

	return isMonotone
}

// Example of runtime file modification
func (b *TimeFileInput) readFile() error {
	if b.file == nil {
		return errors.New("error: TimeFile nil")
	}
	b.file.Seek(0, 0) // also seeks to 0 when piping stdin

	scanner := bufio.NewScanner(b.file) // use default Split: ScanLines

	logger.WriteEvent("reading TimeFileInput file ", b.fname)

	tInterpNew := []float64{}
	xInterpNew := [][]float64{}

	for scanner.Scan() {
		line := scanner.Text()
		words := strings.Fields(line)

		// parse the time
		tInterp, _ := strconv.ParseFloat(words[0], 64)
		tInterpNew = append(tInterpNew, tInterp)

		// parse the values
		xInterp := []float64{}
		for _, word := range words[1:] {
			x, _ := strconv.ParseFloat(word, 64)
			xInterp = append(xInterp, x)
		}
		xInterpNew = append(xInterpNew, xInterp)
	}

	// change the block only if the values are monotone
	if len(tInterpNew) == 0 {
		return errors.New("file NOK")
	} else if checkMonotonicity(tInterpNew) {
		b.tInterp = tInterpNew
		b.xInterp = xInterpNew
		info, _ := b.file.Stat()
		b.modTime = info.ModTime()
	} else {
		return errors.New("warning: invalid TimeFileInput, ignoring")
	}

	return nil
}

func (b *TimeFileInput) refreshFile() {
	// reread the time file
	var openErr error
	b.file, openErr = reopenFile(b.fname, b.file)
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

func (b *TimeFileInput) findInterpSlice() (int, int, float64) {
	// find the lower index
	t := time.Now().Sub(b.start).Seconds()

	if b.cycle {
		tPeriod := b.tInterp[len(b.tInterp)-1]
		if t > tPeriod {
			t = t - tPeriod*float64(int(t/tPeriod))
		}

		if t > tPeriod {
			log.Fatal("bad calc in findInterpSlice()")
		}
	}

	var iLower int
	for iLower = 0; iLower < len(b.tInterp)-1; iLower++ {
		if b.tInterp[iLower] <= t && t < b.tInterp[iLower+1] {
			break
		}
	}

	// upper index
	iUpper := iLower + 1
	if iLower == len(b.tInterp)-1 {
		iUpper = iLower
	}

	// interpolation factor
	dt := b.tInterp[iUpper] - b.tInterp[iLower]
	var alpha float64
	if dt != 0.0 {
		alpha = (t - b.tInterp[iLower]) / dt
	} else {
		alpha = 1.0
	}

	return iLower, iUpper, alpha
}

func (b *TimeFileInput) Update() {
	b.refreshFile()

	// find the time index
	iLower, iUpper, alpha := b.findInterpSlice()

	// interpolate
	x0 := b.xInterp[iLower]
	x1 := b.xInterp[iUpper]

	if len(x0) < len(x1) {
		x0 = append(x0, make([]float64, len(x1)-len(x0))...)
	} else if len(x1) < len(x0) {
		x1 = append(x1, make([]float64, len(x0)-len(x1))...)
	}

	b.out = make([]float64, len(x0))

	if !b.discrete {
		for i := 0; i < len(b.out); i++ { // in place modification
			b.out[i] = (1.0-alpha)*x0[i] + alpha*x1[i]
		}
	} else if alpha < 0.5 {
		b.out = x0
	} else {
		b.out = x1
	}

	b.in = b.out
}

func reopenFile(fname string, file *os.File) (*os.File, error) {
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

func TimeFileInputConstructor(name string, words []string) Block {
	var file *os.File
	fname := words[0]
	var err error
	file, err = reopenFile(fname, file)
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

	b := &TimeFileInput{fname: fname, file: file, start: time.Now(), cycle: cycle, discrete: discrete}

	readErr := b.readFile()
	if readErr != nil {
		logger.WriteFatal("error in TimeFileInputConstructor: "+fname+" not in valid format", readErr)
	}

	return b
}

var TimeFileInputConstructorOk = AddConstructor("TimeFileInput", TimeFileInputConstructor)
