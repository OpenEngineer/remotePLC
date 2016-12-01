package lines

import (
	"../blocks/"
	"log"
	"regexp"
)

func RegexpForkLineConstructor(words []string) Line {
	re1, err1 := regexp.Compile(words[1])
	if err1 != nil {
		log.Fatal("in RegexpForkLineConstructor(), \"", words[1], "\", ", err1)
	}

	b0 := blocks.checkName(words[0])
	b1 := []string{}

	// collect all the matched blocks
	for k, _ := range blocks.Blocks {
		if re1.MatchString(k) {
			b1 = append(b1, k)
		}
	}

	l := &ForkLine{b0: []string{b0}, b1: b1}

	return l
}

var RegexpForkLineConstructorOk = AddConstructor("RegexpForkLine", RegexpForkLineConstructor)
