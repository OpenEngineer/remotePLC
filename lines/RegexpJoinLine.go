package lines

import (
	"../blocks/"
)

func RegexpJoinLineConstructor(name string, words []string, b map[string]blocks.Block) Line {
	b0, _ := getRegexpBlocks(b, words[0])
	b1 := getBlock(b, words[1])

	l := &JoinLine{
		LineData{
			b0:        b0,
			b1:        []blocks.Block{b1},
			DebugName: name,
		},
	}

	return l
}

var RegexpJoinLineConstructorOk = AddConstructor("RegexpJoinLine", RegexpJoinLineConstructor)
