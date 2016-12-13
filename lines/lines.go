package lines

import (
	"../blocks/"
	"../logger/"
	"errors"
	"fmt"
	"log"
	"regexp"
)

type LineData struct {
	b0 []blocks.Block
	b1 []blocks.Block

	n0        []int
	DebugName string
}

type Line interface {
	Transfer() // safe transfer
	Count() (int, int)
	check() bool
	//transfer() // unsafe tranfer (unexported names are not possible with subclasses)
	Info()
}

func (l *LineData) Transfer() {
	log.Fatal("transfer func must be implemented in subclass, " + l.DebugName)
}

func (l *LineData) Info() {
	fmt.Println(l.DebugName, ": ", l.n0, "(", l.b0, ")")
}

// count the number of input float64s
func (l *LineData) Count() (int, int) {
	sumOld := 0
	for _, v := range l.n0 {
		sumOld += v
	}

	newLen := make([]int, len(l.b0))

	sumNew := 0
	for i := 0; i < len(l.b0); i++ {
		newLen[i] = len(l.b0[i].Get())
		sumNew += newLen[i]
	}

	l.n0 = newLen

	return sumNew, sumNew - sumOld
}

func (l *LineData) check() bool {
	ok := true
	var expected int
	var actual int
	if len(l.n0) != len(l.b0) {
		errString := fmt.Sprintln("num input blocks is wrong in Line, ", len(l.n0), " vs ", len(l.b0), " in "+l.DebugName)
		logger.WriteError("LineData.check()", errors.New(errString))
	}

	for i, b := range l.b0 {
		if l.n0[i] != len(b.Get()) {
			ok = false
			expected = l.n0[i]
			actual = len(b.Get())
			break
		}
	}

	if !ok {
		logger.WriteEvent(l.DebugName+", ignoring (numInputs not ok: ", expected, " vs ", actual, ") (ignore this message during init phase)")
	}

	return ok
}

var Constructors = make(map[string](func(string, []string, map[string]blocks.Block) Line))

func AddConstructor(key string, fn func(string, []string, map[string]blocks.Block) Line) bool {
	Constructors[key] = fn
	return true
}

func Construct(name string, constructorType string, args []string,
	b map[string]blocks.Block) Line {
	var l Line
	if constructor, ok := Constructors[constructorType]; ok {
		l = constructor(name, args, b)
	} else {
		log.Fatal("invalid line constructor: ", constructorType)
	}
	return l
}

func ConstructAll(wordsTable [][]string, b map[string]blocks.Block) []Line {
	lines := []Line{}

	for _, words := range wordsTable {
		lines = append(lines, Construct(words[0], words[1], words[2:], b))
		logger.WriteEvent("constructed line: ", words[:])
	}

	return lines
}

func getBlock(bMap map[string]blocks.Block, name string) blocks.Block {
	b, ok := bMap[name]

	if !ok {
		log.Fatal("couldn't find block ", name)
	}

	return b
}

func getBlocks(bMap map[string]blocks.Block, names []string) (bs []blocks.Block) {
	for _, name := range names {
		b := getBlock(bMap, name)
		bs = append(bs, b)
	}

	return
}

func getRegexpBlocks(bMap map[string]blocks.Block, reStr string) (bs []blocks.Block, names []string) {
	re, err := regexp.Compile(reStr)
	logger.WriteFatal("getRegexpBlocks()", err)
	for name, b := range bMap {
		if re.MatchString(name) {
			bs = append(bs, b)
			names = append(names, name)
		}
	}

	return
}
func getDebugName(lineType string, arguments []string) string {
	debugName := lineType

	for _, a := range arguments {
		debugName = debugName + "_" + a
	}

	return debugName
}
