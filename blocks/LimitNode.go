package blocks

import (
  "../logger/"
  "errors"
  "strconv"
)

type LimitNode struct {
	BlockData
  x0 float64
  x1 float64
}

func (b *LimitNode) Put(x []float64) {
	b.in = x
  if len(b.out) != len(b.in) {
    b.out = make([]float64, len(b.in))
  }

  for i, v := range b.in {
    if v < b.x0 {
      b.out[i] = b.x0
    } else if v > b.x1 {
      b.out[i] = b.x1
    } else {
      b.out[i] = v
    }
  }
}

func LimitNodeConstructor(words []string) Block {
	x0, err0 := strconv.ParseFloat(words[0], 64)
  logger.WriteError("in LimitNodeConstructor()", err0)

	x1, err1 := strconv.ParseFloat(words[1], 64)
  logger.WriteError("in LimitNodeConstructor()", err1)

  if x0 > x1 {
    logger.WriteError("in LimitNodeConstructor()", 
      errors.New("x0 needs to be smaller than x1"))
  }

	b := &LimitNode{
    x0: x0,
    x1: x1,
  }
	return b
}

var LimitNodeConstructorOk = AddConstructor("LimitNode", LimitNodeConstructor)
