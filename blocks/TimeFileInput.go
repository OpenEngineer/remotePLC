package blocks

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"time"
)

type TimeFileInput struct {
	InputBlockData
	file    *os.File // usefull for seeking
	scanner *bufio.Scanner
	modTime time.Time
	start   time.Time
	tInterp []float64
	xInterp [][]float64
}

// TODO: log error messages
func checkMonotonicity(tInterp []float64) bool {
	isMonotone := true

	if tInterp[0] != 0.0 {
		isMonotone = false
	}

	for i, t := range tInterp[1:] {
		if t <= tInterp[i] {
			isMonotone = false
		}
	}

	return isMonotone
}

// Example of runtime file modification
func (b *TimeFileInput) readFile() {
	b.file.Seek(0, 0) // also seeks to 0 when piping stdin

	tInterpNew := []float64{}
	xInterpNew := [][]float64{}
	for b.scanner.Scan() {
		tokens := strings.Split(b.scanner.Text(), " ")

		// parse the time
		tInterp, _ := strconv.ParseFloat(tokens[0], 64)
		tInterpNew = append(tInterpNew, tInterp)

		// parse the values
		xInterp := []float64{}
		for _, token := range tokens[1:] {
			x, _ := strconv.ParseFloat(token, 64)
			xInterp = append(xInterp, x)
		}
		xInterpNew = append(xInterpNew, xInterp)
	}

	if checkMonotonicity(tInterpNew) {
		b.tInterp = tInterpNew
		b.xInterp = xInterpNew
		info, _ := b.file.Stat()
		b.modTime = info.ModTime()
	} else {
	}
}

func (b *TimeFileInput) findInterpSlice() (int, int, float64) {
	// find the lower index
	t := time.Now().Sub(b.start).Seconds()
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
	// reread the time file if it has been modified
	info, _ := b.file.Stat()
	if !b.modTime.Equal(info.ModTime()) {
		b.readFile()
	}

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

	for i := 0; i < len(b.out); i++ { // in place modification
		b.out[i] = (1.0-alpha)*x0[i] + alpha*x1[i]
	}

	b.in = b.out
}

func TimeFileInputConstructor(words []string) Block {
	var file *os.File
	if words[0] == "stdin" {
		file = os.Stdin
	} else {
		file, _ = os.Open(words[0])
	}

	scanner := bufio.NewScanner(file) // use default Split: ScanLines
	start := time.Now()
	b := &TimeFileInput{file: file, scanner: scanner, start: start}

	b.readFile()

	return b
}

var TimeFileInputConstructorOk = AddConstructor("TimeFileInput", TimeFileInputConstructor)
