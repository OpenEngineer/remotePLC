package main

import (
	"./graph/"
	"./logger/"
	"errors"
	"os"
	"time"
)

func controlLoop(g *graph.Graph, timeStep time.Duration, saveInterval int,
	logRegexp string) {
	logger.WriteEvent("starting input loop...")
	g.CycleInfinite([]string{"inputs"}, timeStep, 10) // same rate as main loop
	logger.WriteEvent("starting input loop ok")

	// Main loop
	ticker := time.NewTicker(timeStep)
	counter := 0
	logger.WriteEvent("entering control loop...")
	for {
		<-ticker.C

		g.CycleLines()

		g.CycleSerial([]string{"logic"})

		g.CycleLines()

		g.CycleParallel([]string{"outputs"})

		g.CycleSerial([]string{"stops"})
		checkStops(g)

		if counter%saveInterval == 0 {
			g.LogData(logRegexp)
			counter = 0
		}
		counter += 1
	}

	// also save the last time
	g.LogData(logRegexp)
}

func checkStops(g *graph.Graph) {
	v := g.CycleValues([]string{"stops"}, -1.0, func(bname string, x, v float64) float64 {
		if v > x {
			x = v
		}

		if v >= 1.0 {
			logger.WriteFatal("checkStops()", errors.New("in "+bname+": divergence detected"))
		}
		return x
	})
	if v == -1.0 {
		os.Exit(0)
	}
}
