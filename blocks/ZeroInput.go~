package blocks

type ZeroInput struct {
	BlockData
}

func (b *ZeroInput) Update() {
	b.out = []float64{0.0}
}

func ZeroInputConstructor(x []string) Block {
	b := &ZeroInput{}
	return b
}

var ZeroInputOk = AddConstructor("ZeroInput", ZeroInputConstructor)
