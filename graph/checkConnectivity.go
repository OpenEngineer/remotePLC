package graph

import (
	"../logger/"
	"errors"
)

func (g *Graph) checkConnectivity(startBlocks, middleBlocks []string) {
	g.initBlocks(startBlocks, middleBlocks)

	for name, group := range g.b {
		for key, block := range group {
			if len(block.Get()) == 0 {
				logger.WriteFatal("Graph.checkConnectivity()",
					errors.New(name+"["+key+"]"+" not connected"))
			}
		}
	}
}
