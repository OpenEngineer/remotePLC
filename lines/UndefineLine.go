package lines

import (
	"../blocks/"
	"../logger/"
	"errors"
)

type UndefineLine struct {
	LineData
}

// b0 and b1 have (assumed) equal length
func (l *UndefineLine) Transfer() {
	if l.check() {
		for i, b := range l.b0 {
			x := b.Get()
			l.b1[i].Put(x)

			// now overwrite the input block via its Get() with Undefine
			xUndefined := make([]float64, len(x))
			for j := 0; j < len(x); j++ {
				xUndefined[j] = UNDEFINED
			}
			b.Put(xUndefined)
		}
	}
}

func UndefineLineConstructor(name string, words []string, b map[string]blocks.Block) Line {
	if len(words)%2 == 1 {
		logger.WriteFatal("UndefineLineConstructor()", errors.New("unpaired lines"))
	}

	b0 := []blocks.Block{}
	b1 := []blocks.Block{}
	for i := 0; i < len(words); i += 2 {
		b0 = append(b0, getBlock(b, words[i]))
		b1 = append(b1, getBlock(b, words[i+1]))
	}

	l := &UndefineLine{
		LineData{
			b0:        b0,
			b1:        b1,
			DebugName: name,
		},
	}
	return l
}

var UndefineLineConstructorOk = AddConstructor("UndefineLine", UndefineLineConstructor)
