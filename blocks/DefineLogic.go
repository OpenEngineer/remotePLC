package blocks

import (
	"../parser/"
)

type DefineLogic struct {
	BlockData
	defState float64
}

func (b *DefineLogic) Update() {
	isDefined := true
	for _, v := range b.in {
		if v == UNDEFINED {
			isDefined = false
		}
	}

	if isDefined {
		b.out = b.in
	} else {
		if len(b.out) != len(b.in) { // set to defined state, but with same length as b.in
			b.out = MakeDefined(len(b.in), b.defState)
		}
	}
}

func DefineLogicConstructor(name string, words []string) Block {
	var defState float64 // default state if b.in is UNDEFINED
	positional := parser.PositionalArgs(&defState)
	parser.ParsePositionalArgs(words, positional)

	b := &DefineLogic{defState: defState}

	return b
}

var DefineLogicConstructorOk = AddConstructor("DefineLogic", DefineLogicConstructor)
