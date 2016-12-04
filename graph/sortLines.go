package graph

import (
	"sort"
)

type sortGraph struct {
	Graph
	sum  []int
	diff []int
}

func (g *sortGraph) Len() int {
	return len(g.l)
}

func (g *sortGraph) sortCase(i int) int {
	var c int
	if g.sum[i] > 0 && g.diff[i] == 0 {
		c = 0
	} else if g.sum[i] > 0 && g.diff[i] == g.sum[i] {
		c = 1
	} else if g.sum[i] > 0 && g.diff[i] < g.sum[i] {
		c = 2
	} else { // g.sum[i] == 0
		c = 3
	}

	return c
}

func (g *sortGraph) Less(i, j int) bool {
	ci := g.sortCase(i)
	cj := g.sortCase(j)

	b := false
	if ci < cj {
		b = true
	}

	return b
}

func (g *sortGraph) Swap(i, j int) {
	tmpLine := g.l[i]
	tmpSum := g.sum[i]
	tmpDiff := g.diff[i]

	g.l[i] = g.l[j]
	g.sum[i] = g.sum[j]
	g.diff[i] = g.diff[j]

	g.l[j] = tmpLine
	g.sum[j] = tmpSum
	g.diff[j] = tmpDiff
}

func (g *sortGraph) CountLineData() {
	for i, _ := range g.l {
		sum, diff := g.l[i].Count()

		g.sum[i] = sum
		g.diff[i] = diff
	}
}

func (g *Graph) sortLines(startBlocks, middleBlocks []string) {
	numLines := len(g.l)

	s := &sortGraph{
		Graph: Graph{
			b: g.b,
			l: g.l,
		},
		sum:  make([]int, numLines),
		diff: make([]int, numLines),
	}

	s.ClearAll()

	s.CycleParallel(startBlocks)

	for i := 0; i < numLines+1; i++ {
		// update the inputs and the logic
		s.CycleSerial(middleBlocks)

		// count
		s.CountLineData()

		// first the lines where sum > 0 and diff == 0
		// then sum > 0 and diff == sum
		// then sum > 0 and diff < sum
		// then sum == 0
		sort.Stable(s)

		// transfer data before repeating order loop
		s.CycleLines()
	}

	g.l = s.l
}
