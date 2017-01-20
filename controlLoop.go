package main

import (
	"./graph/"
	"./logger/"
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
	stopValue := 0.0
	for stopValue > -1.0 && stopValue < 1.0 {
		<-ticker.C

		g.CycleLines()

		g.CycleSerial([]string{"logic"})

		g.CycleLines()

		g.CycleParallel([]string{"outputs"})

		g.CycleSerial([]string{"stops"})
		stopValue = checkStops(g)

		if counter%saveInterval == 0 {
			g.LogData(logRegexp)
			counter = 0
		}
		counter += 1
	}

	// also save the last time
	g.LogData(logRegexp)

	if stopValue > -1.0 {
		os.Exit(2)
	}
}

func checkStops(g *graph.Graph) float64 {
	// default value for this variable is 0.0
	// this means that without stop criteria the program will run forever
	v := g.CycleValues([]string{"stops"}, 0.0, func(bname string, x, v float64) float64 {
		if v > x {
			x = v
		}

		if v >= 1.0 {
			logger.WriteEvent("checkStops(), ", "in "+bname+": divergence detected")
		}
		return x
	})

	return v
}
