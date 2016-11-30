package blocks

import (
  "../logger/"
  //"fmt"
	"strconv"
)

type ScaleInput struct {
	InputBlockData
	scale  float64
	offset float64
	child  Block
}

func (b *ScaleInput) Update() {
	b.child.Update()
	in := b.child.Get() // To make sure that they are the same size

	if len(b.out) != len(in) {
		b.out = make([]float64, len(in))
	}

	for i, v := range in {
		b.out[i] = b.scale*v + b.offset
	}

	b.in = b.out
  //fmt.Println("ScaleInput", b.out)
}

func ScaleInputConstructor(name string, words []string) Block {
	scale, errScale := strconv.ParseFloat(words[0], 64)
  logger.WriteError("in ScaleInputConstructor()", errScale)

	offset, errOffset := strconv.ParseFloat(words[1], 64)
  logger.WriteError("in ScaleInputConstructor()", errOffset)

	child := Construct(name+"_child", words[2], words[3:])

	b := &ScaleInput{scale: scale, offset: offset, child: child}
	return b
}

var ScaleInputConstructorOk = AddConstructor("ScaleInput", ScaleInputConstructor)
