package blocks

import (
	"../logger/"
	"errors"
	"strconv"
)

type IfElseElseNode struct {
	BlockData
	conditions []float64
	outputs    []float64 // len(outputs) = len(conditions) + 1
}

func (b *IfElseElseNode) eval(x float64) float64 {
	output := b.outputs[0]

	for i, v := range b.conditions {
		if x <= v {
			break
		} else {
			output = b.outputs[i+1]
		}
	}

	return output
}

func (b *IfElseElseNode) Put(x []float64) {
	b.in = x

	tmp := make([]float64, len(x))

	for i, v := range x {
		if v == UNDEFINED {
			tmp[i] = UNDEFINED
		} else {
			tmp[i] = b.eval(x[i])
		}
	}

	b.out = tmp
}

func IfElseElseNodeConstructor(name string, words []string) Block {
	var conditions []float64
	var outputs []float64

	for i := 0; i < len(words); i++ {
		v, err := strconv.ParseFloat(words[i], 64)
		if err != nil {
			logger.WriteError("IfElseElseNodeConstructor()", err)
		}

		if (i % 2) == 0 { // is an output
			outputs = append(outputs, v)
		} else { // is a condition
			conditions = append(conditions, v)
		}
	}

	if !(len(outputs) == len(conditions)+1) {
		logger.WriteError("IfElseElseNodeConstructor()", errors.New("bad number of outputs or conditions"))
	} else if len(outputs) == 0 {
		logger.WriteError("IfElseElseNodeConstructor()", errors.New("must specify at least one output"))
	} else if len(conditions) > 1 { // check that conditions are monotone and unique
		for i := 1; i < len(conditions); i++ {
			if conditions[i] <= conditions[i-1] {
				logger.WriteError("IfElseElseNodeConstructor()", errors.New("conditions must be monotone and unique"))
			}
		}
	}

	b := &IfElseElseNode{
		conditions: conditions,
		outputs:    outputs,
	}

	return b
}

var IfElseElseNodeConstructorOk = AddConstructor("IfElseElseNode", IfElseElseNodeConstructor)
