package blocks

import (
	"../logger/"
	"log"
	"regexp"
	"sort"
)

type BlockModeType int

const (
	REGULAR BlockModeType = iota
	CONNECTIVITY
	STRICT
)

const HIDDEN_SUFFIX_CHAR = "_"

var BlockMode BlockModeType = REGULAR

type Block interface {
	Get() []float64
	Put([]float64)
	Update()
}

type BlockData struct {
	in  []float64
	out []float64
}

func (b *BlockData) Update() {
}

func (b *BlockData) Get() []float64 {
	return b.out
}

func (b *BlockData) Put(x []float64) {
	b.in = x
}

type InputBlockData struct {
	BlockData
}

func (b *InputBlockData) Put(x []float64) {
	b.in = x
	b.out = x
}

type OutputBlockData struct {
	BlockData
}

func (b *OutputBlockData) Get() []float64 {
	b.out = b.in
	return b.out
}

var Blocks = make(map[string]Block)

func checkName(name string) string {
	// if the name doesn't end with underscore then append it, else remove it
	var altName string
	if name[len(name):] == HIDDEN_SUFFIX_CHAR {
		altName = name[:len(name)-1]
	} else {
		altName = name + HIDDEN_SUFFIX_CHAR
	}

	_, ok := Blocks[name]
	_, okAlt := Blocks[altName]

	var checkedName string
	if !ok && okAlt {
		checkedName = altName
	} else if !okAlt && ok {
		checkedName = name
	} else {
		log.Fatal("couldn't find block ", name, " (or ", altName, ")")
	}

	return checkedName
}

func checkNames(names []string) (checkedNames []string) {
	for _, name := range names {
		checkedName := checkName(name)
		checkedNames = append(checkedNames, checkedName)
	}

	return
}

var Constructors = make(map[string]func(string, []string) Block)

func AddConstructor(key string, fn func(string, []string) Block) bool {
	Constructors[key] = fn
	return true
}

func Construct(name string, constructorType string, words []string) Block {
	var b Block
	if constructor, ok := Constructors[constructorType]; ok {
		b = constructor(name, words)
	} else {
		log.Fatal("invalid block constructor: ", constructorType)
	}
	return b
}

func ConstructGlobal(key string, words []string) Block {
	defer logger.WriteEvent("constructed block: ", key, words)
	b := Construct(key, words[0], words[1:])
	Blocks[key] = b
	return b
}

func ConstructAll(wordsTable [][]string) map[string]Block {
	m := make(map[string]Block)
	for _, words := range wordsTable {
		m[words[0]] = ConstructGlobal(words[0], words[1:])
	}
	return m
}

func getSortedNames() (names []string) {
	for name, _ := range Blocks { // eg. inputs, outputs
		names = append(names, name)
	}
	sort.Strings(names)

	return
}

func GetVisibleFields(visibleNameString string) (fields []string, data [][]float64) {
	// names of all the blocks
	names := getSortedNames()

	// visible rule
	visibleName := regexp.MustCompile(visibleNameString)

	for _, name := range names {
		if visibleName.MatchString(name) {
			fields = append(fields, name)
			data = append(data, Blocks[name].Get())
		}
	}

	return
}

func LogData() {
	fields, data := GetVisibleFields(".*[^" + HIDDEN_SUFFIX_CHAR + "]$")

	logger.WriteData(fields, data)
}
