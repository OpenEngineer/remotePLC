package blocks

import "log"
import "regexp"

func RegexpForkLineConstructor(name string, words []string) Block {
	re1, err1 := regexp.Compile(words[1])
	if err1 != nil {
		log.Fatal("in RegexpForkLineConstructor(), \"", words[1], "\", ", err1)
	}

	b0 := words[0]
	b1 := []string{}

	// collect all the matched blocks
	for k, _ := range Blocks {
		if re1.MatchString(k) {
			b1 = append(b1, k)
		}
	}

	b := &ForkLine{b0: b0, b1: b1}

	return b
}

var RegexpForkLineConstructorOk = AddConstructor("RegexpForkLine", RegexpForkLineConstructor)
