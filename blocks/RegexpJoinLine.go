package blocks

import "log"
import "regexp"

func RegexpJoinLineConstructor(name string, words []string) Block {
	re0, err0 := regexp.Compile(words[0])
	if err0 != nil {
		log.Fatal("in RegexpJoinLineConstructor(), \"", words[0], "\", ", err0)
	}

	b0 := []string{}
	b1 := checkName(words[1])

	// collect all the matched blocks
	for k, _ := range Blocks {
		if re0.MatchString(k) {
			b0 = append(b0, k)
		}
	}

	b := &JoinLine{b0: b0, b1: b1}

	return b
}

var RegexpJoinLineConstructorOk = AddConstructor("RegexpJoinLine", RegexpJoinLineConstructor)
