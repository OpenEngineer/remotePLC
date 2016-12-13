package graph

import (
	"../logger/"
)

func (g *Graph) initBlocks(startBlocks, middleBlocks []string) {
	logger.WriteEvent("     ClearAll()...")
	g.ClearAll()
	logger.WriteEvent("     ClearAll() OK")

	logger.WriteEvent("     Cycling inputs()...")
	g.CycleParallel(startBlocks)
	logger.WriteEvent("     Cycling inputs() OK")

	for i := 0; i < len(g.l)+1; i++ {
		g.CycleSerial(middleBlocks)

		g.CycleLines()
	}
}
