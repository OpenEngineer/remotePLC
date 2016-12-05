package graph

import (
	"../blocks/"
	"../lines/"
	"../logger/"
	"fmt"
	"regexp"
	"sort"
	"time"
)

// graph object

type Graph struct {
	b map[string](map[string]blocks.Block)
	l []lines.Line
}

func ConstructGraph(blockTable map[string]([][]string), lineTable [][]string,
	startBlocks, middleBlocks, endBlocks []string) *Graph {

	logger.WriteEvent("constructing graph...")
	g := &Graph{
		b: make(map[string]map[string]blocks.Block),
	}

	logger.WriteEvent("constructing blocks...")
	for k, table := range blockTable {
		g.b[k] = blocks.ConstructAll(table)
	}

	g.l = lines.ConstructAll(lineTable, g.ungroupedBlocks())

	logger.WriteEvent("sorting lines...")
	g.LineInfo()
	g.sortLines(startBlocks, middleBlocks)
	logger.WriteEvent("sorted lines as: ")
	g.LineInfo()

	logger.WriteEvent("checking connectivity...")
	g.checkConnectivity(startBlocks, middleBlocks, endBlocks)
	logger.WriteEvent("connectivity ok")

	logger.WriteEvent("initializing block...")
	g.initBlocks(startBlocks, middleBlocks)
	logger.WriteEvent("blocks initialized")

	logger.WriteEvent("final count...")
	g.CountLineData()
	logger.WriteEvent("final count ok")

	return g
}

func (g *Graph) CycleSerial(names []string) {
	for _, name := range names {
		for _, block := range g.b[name] {
			block.Update()
		}
	}
}

func (g *Graph) CycleParallel(names []string) {
	// count
	count := 0
	for _, name := range names {
		count += len(g.b[name])
	}
	ch := make(chan int, count)

	// fork
	for _, name := range names {
		for _, v := range g.b[name] {
			go func(b blocks.Block) {
				b.Update()
				ch <- 0
			}(v)
		}
	}

	// join
	for _, name := range names {
		for _, _ = range g.b[name] {
			<-ch
		}
	}
}

func (g *Graph) CycleInfinite(names []string, period time.Duration, desync time.Duration) {
	for _, name := range names {
		for _, v := range g.b[name] {
			time.Sleep(desync * time.Millisecond)
			go func(b blocks.Block) {
				ticker := time.NewTicker(period)
				for {
					<-ticker.C
					b.Update()
				}
			}(v)
		}
	}
}

func (g *Graph) CycleLines() {
	for _, l := range g.l {
		l.Transfer()
	}
}

func (g *Graph) ClearAll() {
	for _, group := range g.b {
		for _, block := range group {
			block.Put([]float64{})
		}
	}
}

func (g *Graph) CountLineData() {
	for _, l := range g.l {
		l.Count()
	}
}

func (g *Graph) CycleValues(names []string, init float64, fn func(string, float64, float64) float64) float64 {
	x := init

	for _, name := range names {
		for bname, block := range g.b[name] {
			for _, v := range block.Get() {
				x = fn(bname, x, v)
			}
		}
	}

	return x
}

func (g *Graph) getVisibleFields(visibleNameString string) (fields []string, data [][]float64) {
	// visible rule
	visibleName := regexp.MustCompile(visibleNameString)

  var gnames []string
  for gname, _ := range g.b {
    gnames = append(gnames, gname)
  }
  sort.Strings(gnames)

	for _, gname := range gnames {
		// names of all the blocks
		var names []string
		for name, _ := range g.b[gname] {
			names = append(names, name)
		}
		sort.Strings(names)

		for _, name := range names {
			if visibleName.MatchString(name) {
				fields = append(fields, name)
				data = append(data, g.b[gname][name].Get())
			}
		}

	}
	return
}

func (g *Graph) LogData(namesRegexp string) {
	fields, data := g.getVisibleFields(namesRegexp)

	logger.WriteData(fields, data)
}

func (g *Graph) ungroupedBlocks() map[string]blocks.Block {
	b := make(map[string]blocks.Block)
	for _, group := range g.b {
		for name, block := range group {
			b[name] = block
		}
	}

	return b
}

func (g *Graph) LineInfo() {
	for _, l := range g.l {
		l.Info()
	}
}

func (g *Graph) BlockInfo(name string) {
	fmt.Println(g.b[name])
}
