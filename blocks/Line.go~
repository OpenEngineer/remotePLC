package blocks

type Line struct {
	BlockData
	b0 string
	b1 string
}

func (b *Line) Update() {
	b.in = Blocks[b.b0].Get()
	b.out = b.in
	Blocks[b.b1].Put(b.out)
}

func ConstructLine(words []string) Block {
	b0 := words[0]
	b1 := words[1]

	b := &Line{b0: b0, b1: b1}
	return b
}

var ConstructLineOk = AddConstructor("Line", ConstructLine)
