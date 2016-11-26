package blocks

import (
	"../logger/"
	"log"
  //"net"
)

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

func assertBlock(k string) {
	if _, ok := Blocks[k]; !ok {
		log.Fatal("couldn't find block ", k)
	}
}

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
	defer logger.WriteEvent("constructed block: ", key, words)
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

// a list of open client-server connections, so that they can be reused between the blocks
// the key is: "net:ip:port"
// TODO: really necessary? better to use global conn vars in relevant block types?
//var connections = make(map[string]net.Conn)
