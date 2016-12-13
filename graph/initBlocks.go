package graph

import (
	"../logger/"
)

func (g *Graph) initBlocks(startBlocks, middleBlocks []string) {
	logger.WriteEvent("     ClearAll()...")
	g.ClearAll()
	logger.WriteEvent("     ClearAll() OK")

	g.CycleParallel(startBlocks)

	for i := 0; i < len(g.l)+1; i++ {
		g.CycleSerial(middleBlocks)

		g.CycleLines()
	}
}
