package blocks

import (
  "../logger/"
  "errors"
  "strconv"
)

type BitifyNode struct {
	BlockData
  cutoff float64
}

const defaultBitifyCutoff float64 = 0.0

// example: input [1 1] => 1*2^0 + 1*2^1 = 3
// example: [0 1 2] => 0*2^0 + 1*2^1 + 1*2^2 = 6
// example: [0.01 1000.0] => 1*2^0 + 1*2^1 = 3
// so: numbers > 0.0 are set to 1, and numbers <= 0.0 are set to 0
// the maximum achievable integer in a float64 is: 2^53 (size of mantissa + 1 implicit significand bit, IEEE 754 double precision number)
func (b *BitifyNode) Put(x []float64) {
	b.in = x

  if len(b.in) > 53 {
    logger.WriteError("BitifyNode.Put()",
      errors.New("too many numbers, truncating"))
    b.in = b.in[0:53]
  }

  b.out = []float64{0.0}

  for i, v := range b.in {
    v0 := 0.0
    if v > b.cutoff {
      v0 = 1.0
    }

    b.out[0] += float64(uint(0) << uint(i))*v0
  }
}

func BitifyNodeConstructor(name string, words []string) Block {
  var cutoff float64
  if len(words) == 1 {
    var errCutoff error
    cutoff, errCutoff = strconv.ParseFloat(words[0], 64)
    logger.WriteError("BitifyNodeConstructor()", errCutoff)
  } else {
    cutoff = defaultBitifyCutoff
  }

	b := &BitifyNode{
    cutoff: cutoff,
  }
	return b
}

var BitifyNodeConstructorOk = AddConstructor("BitifyNode", BitifyNodeConstructor)
