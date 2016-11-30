package blocks

import (
  "../logger/"
  "bufio"
	"fmt"
  "os"
	"strconv"
  "strings"
	"time"
)

// TODO: smoother restart
// unit of time is [sec]
type LimitPIDLogic struct {
	BlockData
	kp    float64
	ki    float64
	kd    float64
  x0    float64
  x1    float64
	ePrev []float64
	eInt  []float64
	tPrev time.Time
  name  string
  fname string
  file *os.File // for saving the state of ePrev and eInt
}

func (b *LimitPIDLogic) Update() {
	numIn := len(b.in)
	numPrev := len(b.ePrev)

	// extend b.ePrev and b.eInt if b.in is longer
	if numIn > numPrev {
		b.ePrev = append(b.ePrev, b.in[numPrev:]...)

		for i := 0; i < numIn-numPrev; i++ {
			b.eInt = append(b.eInt, 0.0)
		}
	} else if numIn < numPrev { // shorten if b.in is shorter
		b.ePrev = b.ePrev[0:numIn]
		b.eInt = b.eInt[0:numIn]
	}

	// to get right size:
	if len(b.out) != numIn {
		b.out = make([]float64, numIn)
	}

	// time step size
	t := time.Now()
	dt := t.Sub(b.tPrev).Seconds()

	// modify all arrays inplace:
	for i, e := range b.in {
		de := e - b.ePrev[i]
    dedt := 0.0
    if dt > 0.0 {
      dedt = de / dt
    }

		b.eInt[i] += dt * e
    x := b.kp*e + b.ki*b.eInt[i] + b.kd*dedt

    // apply limiters
    if (x > b.x1) {
      b.out[i] = b.x1
      if b.ki > 1e-8 { 
        b.eInt[i] = (b.x1 - b.kp*e - b.kd*dedt)/b.ki
      } else {
        b.eInt[i] = 0.0
      }
    } else if x < b.x0 {
      b.out[i] = b.x0
      if b.ki > 1e-8 {
        b.eInt[i] = (b.x0 - b.kp*e - b.kd*dedt)/b.ki
      } else {
        b.eInt[i] = 0.0
      }
    } else {
      b.out[i] = x
    }

		b.ePrev[i] = e
	}

	b.tPrev = t
  b.saveState()
}

func (b *LimitPIDLogic) saveState() {
  b.file.Seek(0,0)

  // First column contains the previous error, second column contains the error integral
  for i, v := range b.ePrev {
    fmt.Fprintln(b.file, v, b.eInt[i])
  }
}

func (b *LimitPIDLogic) loadState() {
  // init
  b.ePrev = []float64{}
  b.eInt = []float64{}


  file, errOpen := os.Open(b.fname)
  if errOpen == nil {
    // read the lines, and within those lines read the columns
    scanner := bufio.NewScanner(file)

    for scanner.Scan() {
      line := scanner.Text()
      words := strings.Fields(line)

      if len(words) == 2 {
        ePrev, errPrev := strconv.ParseFloat(words[0], 64)
        eInt, errInt := strconv.ParseFloat(words[1], 64)

        if errPrev == nil && errInt == nil {
          b.ePrev = append(b.ePrev, ePrev)
          b.eInt = append(b.eInt, eInt)
          logger.WriteEvent("LimitPIDLogic.loadState(), :", ePrev, eInt)
        } else {
          logger.WriteError("LimitPIDLogic.loadState()", errPrev)
          logger.WriteError("LimitPIDLogic.loadState()", errInt)
        }
      }
    }
  } else {
    logger.WriteEvent("couldn't restart pid from a state file, restarting from scratch")
  }
}

func LimitPIDLogicConstructor(name string, words []string) Block {
	kp, _ := strconv.ParseFloat(words[0], 64)
	ki, _ := strconv.ParseFloat(words[1], 64)
	kd, _ := strconv.ParseFloat(words[2], 64)
	x0, _ := strconv.ParseFloat(words[3], 64)
	x1, _ := strconv.ParseFloat(words[4], 64)

	tPrev := time.Now()

  fname := name+"_state.dat"

  b := &LimitPIDLogic{kp: kp, ki: ki, kd: kd, x0: x0, x1: x1, tPrev: tPrev, name: name, fname: fname}

  // only open file for reading:
  b.loadState()

  // anad finally load the file for writing
  var errCreate error
  b.file, errCreate = os.Create(fname)
  logger.WriteError("LimitPIDLogicConstructor()", errCreate)
  
	return b
}

var LimitPIDLogicConstructorOk = AddConstructor("LimitPIDLogic", LimitPIDLogicConstructor)
