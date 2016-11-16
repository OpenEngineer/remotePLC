package tables

type ZeroInput struct{}

func (z ZeroInput) Get() float64 {
	return 0.0
}

func ZeroInputConstructor(x []string) Input {
	return ZeroInput{}
}

var ZeroInputOk = AddInputConstructor("ZeroInput", ZeroInputConstructor)
