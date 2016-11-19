package blocks

import "time"
import "log"

type TimeOutStop struct {
	BlockData
	seconds float64
	start   time.Time
}

func (b *TimeOutStop) Update() {
	d := time.Since(b.start)

	t := d.Seconds()
	if t > b.seconds {
		b.out = []float64{1.0 + t/b.seconds}
	} else {
		b.out = []float64{-1.0 - t/b.seconds}
	}
}

func TimeOutStopConstructor(words []string) Block {
	d, err := time.ParseDuration(words[0])
	if err != nil {
		log.Fatal("in TimeOutStopConstructor, \"", words[0], "\", ", err)
	}
	seconds := d.Seconds()
	start := time.Now()

	b := &TimeOutStop{seconds: seconds, start: start}
	return b
}

var TimeOutStopConstructorOk = AddConstructor("TimeOutStop", TimeOutStopConstructor)
