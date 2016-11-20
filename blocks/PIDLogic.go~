package blocks

import (
	"strconv"
	"time"
)

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
	// extend b.ePrev and b.eInt if b.in is longer
	if len(b.in) > len(b.ePrev) {
		b.ePrev = append(b.ePrev, b.in[len(b.ePrev):len(b.in)]...)

		for i := 0; i < len(b.in)-len(b.ePrev); i++ {
			b.eInt = append(b.eInt, 0.0)
		}
	} else if len(b.in) < len(b.ePrev) { // shorten if b.in is shorter
		b.ePrev = b.ePrev[0:len(b.in)]
		b.eInt = b.eInt[0:len(b.in)]
	}

	// to get right size:
	b.out = b.in

	// time step size
	t := time.Now()
	dt := t.Sub(b.tPrev).Seconds()

	// modify all arrays inplace:
	for i, e := range b.in {
		de := e - b.ePrev[i]
		b.eInt[i] += dt * de
		dedt := de / dt
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
