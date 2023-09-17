package main

import (
	"fmt"
	"hw/hm"

	"github.com/bradhe/stopwatch"
)

func main() {
	watch := stopwatch.Start()
	var generator = hm.HeightMapGenerator{}
	var heightMap = generator.Generate(0.5, 12)
	heightMap.Print()
	watch.Stop()
	fmt.Printf("Milliseconds elapsed: %v\n", watch.Seconds())
}
