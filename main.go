package main

import (
	"hw/hm"
)

func main() {
	var generator = hm.HeightMapGenerator{}
	var heightMap = generator.Generate(3, 0.5)
	heightMap.Print()
}
