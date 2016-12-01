package lines

import (
	"../blocks/"
	"log"
)

// multiple lines are possible
type Line struct {
	LineData
}

// b0 and b1 have (assumed) equal length
func (l *Line) Transfer() {
	for i, v := range l.b0 {
		x := blocks.Blocks[v].Get()
		blocks.Blocks[l.b1[i]].Put(x)
	}
}

func LineConstructor(words []string) Line {
	if len(words)%2 == 1 {
		log.Fatal("in LineConstructor, ", words, ",error: unpaired lines")
	}

	b0 := []string{}
	b1 := []string{}
	for i := 0; i < len(words); i += 2 {
		b0 = append(b0, blocks.checkName(words[i]))
		b1 = append(b1, blocks.checkName(words[i+1]))
	}

	l := &Line{b0: b0, b1: b1}
	return l
}

var LineConstructorOk = AddConstructor("Line", LineConstructor)
