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

const UNDEFINED float64 = 1e300

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

func SafeCopy(numCopy int, src []float64, numDst int) []float64 {
	dst := make([]float64, numDst)

	// copy numCopy floats from src to dst
	for i := 0; i < numCopy; i++ {
		if i < len(src) {
			dst[i] = src[i]
		} else {
			if i < len(dst) {
				dst[i] = 0.0
			} else {
				dst = append(dst, 0.0)
			}
		}
	}
	return dst
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
