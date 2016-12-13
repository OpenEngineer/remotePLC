package graph

import (
	"../logger/"
	"errors"
)

func (g *Graph) checkConnectivity(startBlocks, middleBlocks, endBlocks []string) {
	logger.WriteEvent("  initBlocks()...")
	g.initBlocks(startBlocks, middleBlocks)

	logger.WriteEvent("  initBlocks() OK")
	for _, l := range g.l {
		sum, _ := l.Count()

		if sum == 0 {
			l.Info()
			logger.WriteFatal("Graph.checkConnectivity()",
				errors.New("unconnected line (could be due to circularity)"))
		}
	}
	for _, gname := range endBlocks {
		for key, block := range g.b[gname] {
			if len(block.Get()) == 0 {
				logger.WriteFatal("Graph.checkConnectivity()",
					errors.New(gname+"["+key+"]"+" not connected"))
			}
		}
	}
}
