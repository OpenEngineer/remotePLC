package blocks

import (
	"../logger/"
	"../parser/"
	"errors"
	"math"
)

type WrapNode struct {
	BlockData
	lower float64
	upper float64

	fn func(x, lower, upper float64) float64
}

func (b *WrapNode) Put(x []float64) {
	b.in = x
	if len(b.out) != len(b.in) {
		b.out = make([]float64, len(b.in))
	}

	for i, v := range b.in {
		b.out[i] = b.fn(v, b.lower, b.upper)
	}
}

func WrapNodeConstructor(name string, words []string) Block {
	var lower float64
	var upper float64
	mode := "bound"
	positional := parser.PositionalArgs(&lower, &upper)
	optional := parser.OptionalArgs("Mode", &mode)

	parser.ParseArgs(words, positional, optional)

	fns := make(map[string]func(float64, float64, float64) float64)
	fns["bound"] = WrapNodeBound
	fns["cycle"] = WrapNodeCycle

	var fn func(float64, float64, float64) float64
	var ok bool
	if fn, ok = fns[mode]; !ok {
		validModes := ""
		for key, _ := range fns {
			validModes = validModes + ", " + key
		}

		logger.WriteFatal("ReductionNodeConstructor()", errors.New("mode "+mode+" not found in \""+validModes+"\""))
		return nil
	}

	if upper < lower {
		logger.WriteError("ReductionNodeConstructor()", errors.New("bad bounds"))
	}

	b := &WrapNode{
		lower: lower,
		upper: upper,
		fn:    fn,
	}
	return b
}

var WrapNodeConstructorOk = AddConstructor("WrapNode", WrapNodeConstructor)

func WrapNodeBound(x, lower, upper float64) float64 {
	xBound := x
	if xBound > upper {
		xBound = upper
	} else if xBound < lower {
		xBound = lower
	}

	return xBound
}

func WrapNodeCycle(x, lower, upper float64) float64 {
	xCycle := x
	if xCycle > upper {
		xCycle = xCycle - (upper-lower)*math.Floor((x-lower)/(upper-lower))
	} else if xCycle < lower {
		xCycle = xCycle + (upper-lower)*math.Floor((upper-x)/(upper-lower))
	}

	return xCycle
}
