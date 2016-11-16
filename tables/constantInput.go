package tables

import "strconv"

type ConstantInput struct {
	a float64
}

func (c ConstantInput) Get() float64 {
	return c.a
}

func ConstructConstantInput(x []string) Input {
	c := ConstantInput{}
	c.a, _ = strconv.ParseFloat(x[0], 64)
	return c
}

var ConstantInputOk = AddInputConstructor("ConstantInput", ConstructConstantInput)
