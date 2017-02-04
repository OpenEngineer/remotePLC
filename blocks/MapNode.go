package blocks

import (
	"../logger/"
	"../parser/"
	"errors"
	//"fmt"
	"math"
	"os"
	"strconv"
	"time"
)

// implemented as copy of TimeFileInput
// TODO: factor out common functionality

type MapNode struct {
	BlockData

	fname   string
	file    *os.File // usefull for seeking
	modTime time.Time

	numInput  int
	numOutput int
	xInterp   [][]float64
	yInterp   [][]float64

	// options
	mode string // defaults to "interpolate", other options are "nearest" and "exact"
}

// Example of runtime file modification
func (b *MapNode) readFile() error {
	tokens := parser.TokenizeFile(b.fname)
	//fmt.Println("parsed file as: ", tokens)

	// temporary storage to check for validity
	xInterpTmp := [][]float64{}
	yInterpTmp := [][]float64{}

	isOk := true
	for _, l := range tokens {
		// row data
		xInterpRow := []float64{}
		yInterpRow := []float64{}
		if len(l) == b.numInput+b.numOutput {
			for i := 0; i < b.numInput; i++ {
				xInterp, err := strconv.ParseFloat(l[i], 64)
				if err != nil {
					isOk = false
					break
				}
				xInterpRow = append(xInterpRow, xInterp)
			}

			if !isOk {
				break
			}

			for i := b.numInput; i < b.numOutput+b.numInput; i++ {
				yInterp, err := strconv.ParseFloat(l[i], 64)
				if err != nil {
					isOk = false
					break
				}
				yInterpRow = append(yInterpRow, yInterp)
			}

			if !isOk {
				break
			}

		} else {
			logger.WriteEvent(b.fname, " contains row not of length ", b.numInput+b.numOutput, " (but ", len(l), ")")
			isOk = false
			break
		}

		xInterpTmp = append(xInterpTmp, xInterpRow)
		yInterpTmp = append(yInterpTmp, yInterpRow)
	}

	if !isOk {
		return errors.New("warning: invalid MapNode, ignoring")
	} else {
		b.xInterp = xInterpTmp
		b.yInterp = yInterpTmp

		info, _ := b.file.Stat()
		b.modTime = info.ModTime()
		return nil
	}
}

// it is assumed that the parsing the file takes some time
// so we check if it has had an update, this requires "reopening" the file though
func (b *MapNode) refreshFile() {
	// reread the time file (i.e. update Stat)
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

// linear interpolation between the two closest bounding numbers (or clipped if no upper or lower bounds found)
func (b *MapNode) Interp1D(x float64) []float64 {
	// find the numbers the are closest and smaller, and closer and bigger
	iLower := -1
	dMinLower := 1e300

	iUpper := -1
	dMinUpper := 1e300

	for i, v := range b.xInterp {
		dLower := x - v[0]
		if dLower >= 0 && dLower < dMinLower {
			dMinLower = dLower
			iLower = i
		}

		dUpper := v[0] - x
		if dUpper > 0 && dUpper < dMinUpper {
			dMinUpper = dUpper
			iUpper = i
		}
	}

	yInterp := make([]float64, b.numOutput)

	if iLower == -1 && iUpper == -1 {
		for i, _ := range yInterp {
			yInterp[i] = UNDEFINED
		}
	} else if iLower > -1 && iUpper == -1 {
		yInterp = b.yInterp[iLower]
	} else if iUpper > -1 && iLower == -1 {
		yInterp = b.yInterp[iUpper]
	} else {
		// interpolation fraction
		f := dMinLower / (dMinUpper + dMinLower)
		for i, _ := range yInterp {
			yInterp[i] = (1.0-f)*b.yInterp[iLower][i] + f*b.yInterp[iUpper][i]
		}
	}

	return yInterp
}

func (b *MapNode) NearestND(x []float64) []float64 {
	iNearest := -1
	dMin := 1e300 // euclidian distance

	// check which vector is closest
	for i, v := range b.xInterp {
		// loop the vector components
		d := 0.0
		for j, vv := range v {
			d = d + (vv-x[j])*(vv-x[j])
		}

		if d < dMin {
			dMin = d
			iNearest = i
		}
	}

	yInterp := make([]float64, b.numOutput)
	if iNearest == -1 {
		for i, _ := range yInterp {
			yInterp[i] = UNDEFINED
		}
	} else {
		yInterp = b.yInterp[iNearest]
	}
	return yInterp
}

// inverse distance weighting
func (b *MapNode) InterpND(x []float64) []float64 {

	// calculate the weights
	w := make([]float64, len(b.xInterp))
	wSum := 0.0
	isExact := false
	iExact := -1
	for i, v := range b.xInterp {
		d := 0.0
		for j, vv := range v {
			d = d + (vv-x[j])*(vv-x[j])
		}

		if d == 0.0 {
			isExact = true
			iExact = i
			break
		} else {
			w[i] = math.Sqrt(1.0 / d)
			wSum = wSum + w[i]
		}
	}

	// calculate the interpolated vector
	yInterp := make([]float64, b.numOutput)
	if isExact {
		yInterp = b.yInterp[iExact]
	} else {
		for i, v := range b.yInterp {
			for j, vv := range v {
				yInterp[j] = yInterp[j] + vv*w[i]/wSum
			}
		}
	}

	return yInterp
}

func (b *MapNode) Exact(x []float64) []float64 {
	isExact := false
	iExact := -1

	for i, v := range b.xInterp {
		isExactLocal := true
		for j, vv := range v {
			if vv != x[j] {
				isExactLocal = false
				//fmt.Println("broke exactness at ", vv, " vs ", x[j], "at pos ", j)
				break
			}
		}

		if isExactLocal {
			isExact = true
			iExact = i
			//fmt.Println("found an exact match at line ", i)
			break
		}
	}

	yInterp := make([]float64, b.numOutput)
	if isExact {
		yInterp = b.yInterp[iExact]
	} else {
		for i, _ := range yInterp {
			yInterp[i] = UNDEFINED
		}
	}
	return yInterp
}

func (b *MapNode) Put(x []float64) {
	if len(x) != b.numInput {
		return // don't do anything
	}
	b.refreshFile()

	// do different things depending on the settings
	var yInterp []float64

	switch mode := b.mode; mode {
	case "exact":
		yInterp = b.Exact(x) // Undefined if not exact
	case "nearest":
		yInterp = b.NearestND(x)
	case "interpolate":
		if b.numInput == 1 {
			yInterp = b.Interp1D(x[0])
		} else {
			yInterp = b.InterpND(x)
		}
	default:
		logger.WriteError("MapNode.Put()", errors.New("MapNode mode not recognized"))
	}

	b.in = x
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
	var fname string
	var numInput int
	var numOutput int
	var mode string
	mode = "interpolate" // or possible values are "nearest" and "exact"

	positional := parser.PositionalArgs(&fname, &numInput, &numOutput)
	optional := parser.OptionalArgs("Mode", &mode)

	parser.ParseArgs(words, positional, optional)

	var file *os.File
	var err error
	file, err = reopenFile_(fname, file)
	logger.WriteError("MapNodeConstructor()", err)

	b := &MapNode{
		fname:     fname,
		file:      file,
		numInput:  numInput,
		numOutput: numOutput,
		mode:      mode,
	}

	readErr := b.readFile()
	logger.WriteError("error in MapNodeConstructor: "+b.fname+" not in valid format", readErr)

	logger.WriteEvent("constructed MapNode " + name)
	return b
}

var MapNodeConstructorOk = AddConstructor("MapNode", MapNodeConstructor)
