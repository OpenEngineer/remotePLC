package blocks

import (
	"time"
)

type TimeInput struct {
	InputBlockData
	start time.Time
}

func (b *TimeInput) Update() {
	b.out = []float64{time.Now().Sub(b.start).Seconds()}
	b.in = b.out
}

func TimeInputConstructor(name string, words []string) Block {
	b := &TimeInput{start: time.Now()}
	return b
}

var TimeInputConstructorOk = AddConstructor("TimeInput", TimeInputConstructor)
