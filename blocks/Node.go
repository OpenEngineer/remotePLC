package blocks

import ()

type Node struct {
	BlockData
}

func (b *Node) Put(x []float64) {
	b.in = x
	b.out = x
}

func NodeConstructor(name string, words []string) Block {
	b := &Node{}
	return b
}

var NodeConstructorOk = AddConstructor("Node", NodeConstructor)
