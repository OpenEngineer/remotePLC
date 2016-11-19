package blocks

import "strconv"

type ScaledInput struct {
	InputBlockData
	scale  float64
	offset float64
	child  Block
}

func (b *ScaledInput) Update() {
	b.child.Update()
	b.out = b.child.Get() // To make sure that they are the same size

	for i, v := range b.out {
		b.out[i] = b.scale*v + b.offset
	}

	b.in = b.out
}

func ScaledInputConstructor(words []string) Block {
	scale, _ := strconv.ParseFloat(words[0], 64)
	offset, _ := strconv.ParseFloat(words[1], 64)
	child := Construct(words[2], words[3:])

	b := &ScaledInput{scale: scale, offset: offset, child: child}
	return b
}

var ScaledInputConstructorOk = AddConstructor("ScaledInput", ScaledInputConstructor)
