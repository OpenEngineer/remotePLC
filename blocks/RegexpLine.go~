package blocks

import "log"
import "regexp"

func RegexpLineConstructor(words []string) Block {
	re0, err0 := regexp.Compile(words[0])
	if err0 != nil {
		log.Fatal("in RegexpLineConstructor(), \"", words[0], "\", ", err0)
	}
	re1, err1 := regexp.Compile(words[1])
	if err1 != nil {
		log.Fatal("in RegexpLineConstructor(), \"", words[1], "\", ", err1)
	}

	b0 := []string{}
	b1 := []string{}

	// collect all the matched blocks
	for k, _ := range Blocks {
		if re0.MatchString(k) {
			b0 = append(b0, k)
		}

		if re1.MatchString(k) {
			b1 = append(b1, k)
		}
	}

	// truncate to minimum of either
	if len(b0) < len(b1) {
		b1 = b1[:len(b0)]
	} else if len(b0) > len(b1) {
		b0 = b0[:len(b1)]
	}

	b := &Line{b0: b0, b1: b1}

	return b
}

var RegexpLineConstructorOk = AddConstructor("RegexpLine", RegexpLineConstructor)
