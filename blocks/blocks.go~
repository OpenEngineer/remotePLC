package blocks

import "log"

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

var Blocks = make(map[string]Block)

var Constructors = make(map[string]func([]string) Block)

func AddConstructor(key string, fn func([]string) Block) bool {
	Constructors[key] = fn
	return true
}

func Construct(x string, y []string) Block {
	var b Block
	if constructor, ok := Constructors[x]; ok {
		b = constructor(y)
	} else {
		log.Fatal("invalid block constructor: ", x)
	}
	return b
}

func ConstructGlobal(key string, words []string) Block {
	b := Construct(words[0], words[1:])
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
