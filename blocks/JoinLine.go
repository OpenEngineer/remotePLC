package blocks

type JoinLine struct {
	BlockData
	b0 []string
	b1 string
}

func (b *JoinLine) Update() {
	b.in = []float64{}

	for _, v := range b.b0 {
		b.in = append(b.in, Blocks[v].Get()...)
	}
	Blocks[b.b1].Put(b.in)

	b.out = b.in
}

func JoinLineConstructor(name string, words []string) Block {
	b0 := words[1:]
	b1 := words[0]

  assertBlock(b1)
  for _, v := range b0 {
    assertBlock(v)
  }
	b := &JoinLine{b0: b0, b1: b1}
	return b
}

var JoinLineConstructorOk = AddConstructor("JoinLine", JoinLineConstructor)
