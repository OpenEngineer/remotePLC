package blocks

import (
	"log"
)

type BlockModeType int

const (
	REGULAR BlockModeType = iota
	CONNECTIVITY
	STRICT
)

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

var Constructors = make(map[string]func(string, []string) Block)

func AddConstructor(key string, fn func(string, []string) Block) bool {
	Constructors[key] = fn
	return true
}

func Construct(name string, constructorType string, args []string) Block {
	var b Block
	if constructor, ok := Constructors[constructorType]; ok {
		b = constructor(name, args)
	} else {
		log.Fatal("invalid block constructor: ", constructorType)
	}
	return b
}

func ConstructAll(wordsTable [][]string) map[string]Block {
	m := make(map[string]Block)
	for _, words := range wordsTable {
		m[words[0]] = Construct(words[0], words[1], words[2:])
	}
	return m
}
