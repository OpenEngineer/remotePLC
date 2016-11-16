package tables

import "strconv"

type ScaledInput struct {
	a     float64
	b     float64
	input Input
}

func (s ScaledInput) Get() float64 {
	return s.input.Get()*s.a + s.b
}

func ConstructScaledInput(x []string) Input {
	s := ScaledInput{}
	s.a, _ = strconv.ParseFloat(x[0], 64)
	s.b, _ = strconv.ParseFloat(x[1], 64)

	s.input = InputConstructors[x[2]](x[3:])
	return s
}

var ScaledInputOk = AddInputConstructor("ScaledInput", ConstructScaledInput)
