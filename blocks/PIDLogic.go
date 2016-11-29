package blocks

import (
	"strconv"
	"time"
	//"fmt"
)

// TODO: smoother restart
// unit of time is [sec]
type PIDLogic struct {
	BlockData
	kp    float64
	ki    float64
	kd    float64
	ePrev []float64
	eInt  []float64
	tPrev time.Time
}

func (b *PIDLogic) Update() {
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
		dedt := de / dt
		b.eInt[i] += dt * e
		b.out[i] = b.kp*b.in[i] + b.ki*b.eInt[i] + b.kd*dedt
		b.ePrev[i] = e
	}

	b.tPrev = t
}

func PIDLogicConstructor(words []string) Block {
	kp, _ := strconv.ParseFloat(words[0], 64)
	ki, _ := strconv.ParseFloat(words[1], 64)
	kd, _ := strconv.ParseFloat(words[2], 64)

	tPrev := time.Now()
	b := &PIDLogic{kp: kp, ki: ki, kd: kd, tPrev: tPrev}
	return b
}

var PIDLogicConstructorOk = AddConstructor("PIDLogic", PIDLogicConstructor)
