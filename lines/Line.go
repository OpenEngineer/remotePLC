package lines

import (
	"../blocks/"
	"log"
)

type SimpleLine struct { // must be different
	LineData
}

// b0 and b1 have (assumed) equal length
func (l *SimpleLine) Transfer() {
	if l.check() {
		for i, b := range l.b0 {
			x := b.Get()
			l.b1[i].Put(x)
		}

	}
}

func LineConstructor(name string, words []string, b map[string]blocks.Block) Line {
	if len(words)%2 == 1 {
		log.Fatal("in LineConstructor, ", words, ",error: unpaired lines")
	}

	b0 := []blocks.Block{}
	b1 := []blocks.Block{}
	for i := 0; i < len(words); i += 2 {
		b0 = append(b0, getBlock(b, words[i]))
		b1 = append(b1, getBlock(b, words[i+1]))
	}

	l := &SimpleLine{
		LineData{
			b0:        b0,
			b1:        b1,
			DebugName: name,
		},
	}
	return l
}

var LineConstructorOk = AddConstructor("Line", LineConstructor)
