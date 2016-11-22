package blocks

import (
	"log"
	"strconv"
)

type ScaledInput struct {
	InputBlockData
	scale  float64
	offset float64
	child  Block
}

func (b *ScaledInput) Update() {
	b.child.Update()
	in := b.child.Get() // To make sure that they are the same size

	if len(b.out) != len(in) {
		b.out = make([]float64, len(in))
	}

	for i, v := range in {
		b.out[i] = b.scale*v + b.offset
	}

	b.in = b.out
}

func ScaledInputConstructor(words []string) Block {
	scale, errScale := strconv.ParseFloat(words[0], 64)
	if errScale != nil {
		log.Fatal("in ScaledInputConstructor(), ", errScale)
	}

	offset, errOffset := strconv.ParseFloat(words[1], 64)
	if errOffset != nil {
		log.Fatal("in ScaledInputConstructor(), ", errOffset)
	}

	child := Construct(words[2], words[3:])

	b := &ScaledInput{scale: scale, offset: offset, child: child}
	return b
}

var ScaledInputConstructorOk = AddConstructor("ScaledInput", ScaledInputConstructor)
