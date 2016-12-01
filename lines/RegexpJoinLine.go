package lines

import (
	"log"
	"regexp"
)

func RegexpJoinLineConstructor(words []string) Line {
	re0, err0 := regexp.Compile(words[0])
	if err0 != nil {
		log.Fatal("in RegexpJoinLineConstructor(), \"", words[0], "\", ", err0)
	}

	b0 := []string{}
	b1 := blocks.checkName(words[1])

	// collect all the matched blocks
	for k, _ := range blocks.Blocks {
		if re0.MatchString(k) {
			b0 = append(b0, k)
		}
	}

	l := &JoinLine{b0: b0, b1: []string{b1}}

	return l
}

var RegexpJoinLineConstructorOk = AddConstructor("RegexpJoinLine", RegexpJoinLineConstructor)
