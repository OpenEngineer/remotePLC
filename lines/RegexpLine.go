package lines

import (
	"../blocks/"
)

func RegexpLineConstructor(name string, words []string, b map[string]blocks.Block) Line {
	b0, _ := getRegexpBlocks(b, words[0])
	b1, _ := getRegexpBlocks(b, words[1])

	// truncate to minimum of either
	if len(b0) < len(b1) {
		b1 = b1[:len(b0)]
	} else if len(b0) > len(b1) {
		b0 = b0[:len(b1)]
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

var RegexpLineConstructorOk = AddConstructor("RegexpLine", RegexpLineConstructor)
