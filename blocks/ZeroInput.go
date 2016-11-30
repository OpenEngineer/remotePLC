package blocks

type ZeroInput struct {
	InputBlockData
}

func (b *ZeroInput) Update() {
	b.out = []float64{0.0}
	b.in = b.out
}

func ZeroInputConstructor(name string, words []string) Block {
	b := &ZeroInput{}
	return b
}

var ZeroInputConstructorOk = AddConstructor("ZeroInput", ZeroInputConstructor)
