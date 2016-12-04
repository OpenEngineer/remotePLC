package graph

func (g *Graph) initBlocks(startBlocks, middleBlocks []string) {
	g.ClearAll()

	g.CycleParallel(startBlocks)

	for i := 0; i < len(g.l)+1; i++ {
		g.CycleSerial(middleBlocks)

		g.CycleLines()
	}
}
