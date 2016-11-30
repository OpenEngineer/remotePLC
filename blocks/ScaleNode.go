package blocks

import (
  "../logger/"
  "strconv"
)

type ScaleNode struct {
	BlockData
  scale float64
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
	scale, errScale := strconv.ParseFloat(words[0], 64)
  logger.WriteError("in ScaleNodeConstructor()", errScale)

	offset, errOffset := strconv.ParseFloat(words[1], 64)
  logger.WriteError("in ScaleNodeConstructor()", errOffset)

	b := &ScaleNode{
    scale: scale,
    offset: offset,
  }
	return b
}

var ScaleNodeConstructorOk = AddConstructor("ScaleNode", ScaleNodeConstructor)
