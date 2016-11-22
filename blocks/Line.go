package blocks

import "log"

// multiple lines are possible
type Line struct {
	BlockData
	b0 []string
	b1 []string
}

// b0 and b1 have equal length
func (b *Line) Update() {
	b.in = []float64{}
	for i, v := range b.b0 {
		x := Blocks[v].Get()
		Blocks[b.b1[i]].Put(x)

		b.in = append(b.in, x...)
	}
	b.out = b.in
}

func LineConstructor(words []string) Block {
	if len(words)%2 == 1 {
		log.Fatal("in LineConstructor, ", words, ",error: unpaired lines")
	}

	b0 := []string{}
	b1 := []string{}
	for i := 0; i < len(words); i += 2 {
		assertBlock(words[i])
		assertBlock(words[i+1])
		b0 = append(b0, words[i])
		b1 = append(b1, words[i+1])
	}

	b := &Line{b0: b0, b1: b1}
	return b
}

var LineConstructorOk = AddConstructor("Line", LineConstructor)
