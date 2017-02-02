package blocks

import (
	"../parser/"
)

type ScaleNode struct {
	BlockData
	scale  float64
	offset float64
}

func (b *ScaleNode) Put(x []float64) {
	b.in = x
	if len(b.out) != len(b.in) {
		b.out = make([]float64, len(b.in))
	}

	for i, v := range b.in {
		b.out[i] = b.scale*v + b.offset
	}
}

func ScaleNodeConstructor(name string, words []string) Block {
	var scale float64
	var offset float64
	positional := parser.PositionalArgs(&scale, &offset)
	parser.ParsePositionalArgs(words, positional)

	b := &ScaleNode{
		scale:  scale,
		offset: offset,
	}
	return b
}

var ScaleNodeConstructorOk = AddConstructor("ScaleNode", ScaleNodeConstructor)
