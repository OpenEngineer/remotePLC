package lines

import (
	"../blocks/"
)

func RegexpForkLineConstructor(words []string, b map[string]blocks.Block) Line {
	b0 := getBlock(b, words[0])
	b1, _ := getRegexpBlocks(b, words[1])

	l := &ForkLine{
		LineData{
			b0:        []blocks.Block{b0},
			b1:        b1,
			DebugName: getDebugName("RegexpForkLine", words),
		},
	}

	return l
}

var RegexpForkLineConstructorOk = AddConstructor("RegexpForkLine", RegexpForkLineConstructor)
