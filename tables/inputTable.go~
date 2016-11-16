package tables

type Input interface {
	Get() float64
}

var Inputs = make(map[string]Input)

var InputConstructors = make(map[string]func([]string) Input)

func AddInputConstructor(key string, fn func([]string) Input) bool {
	InputConstructors[key] = fn
	return true
}

func ConstructInputs(x [][]string) {
	for _, v := range x {
		Inputs[v[0]] = InputConstructors[v[1]](v[2:])
	}
}
