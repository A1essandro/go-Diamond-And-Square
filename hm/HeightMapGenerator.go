package hm

import (
	"fmt"
	"math"
	"math/rand"
	"sync"
)

// Generator of Height Map
type HeightMapGenerator struct {
}

// Internal struct for contextions
type generateContext struct {
	Size        int
	Persistance float32
	HeightMap   *HeightMap
}

// Type of height map
type HeightMap [][]float32

// Main method of HeightMapGenerator
// Generates HeightMap by contextions
func (g HeightMapGenerator) Generate(size int, persistance float32) HeightMap {
	realSize := int(math.Pow(2, float64(size))) + 1
	heightMap := make(HeightMap, realSize)
	for i := range heightMap {
		(heightMap)[i] = make([]float32, realSize)
	}

	context := generateContext{
		Size:        realSize,
		Persistance: persistance,
		HeightMap:   &heightMap,
	}

	g.initCorners(context)
	g.divide(context, realSize)

	return heightMap
}

// Initializing corners of height map
func (g HeightMapGenerator) initCorners(context generateContext) {
	lastIndex := context.Size - 1
	hm := *context.HeightMap

	hm[0][0] = g.getOffset(context.Size, context)
	hm[0][lastIndex] = g.getOffset(context.Size, context)
	hm[lastIndex][0] = g.getOffset(context.Size, context)
	hm[lastIndex][lastIndex] = g.getOffset(context.Size, context)
}

// Getting random offset of height
func (g HeightMapGenerator) getOffset(stepSize int, context generateContext) float32 {
	return float32(stepSize) / float32(context.Size) * rand.Float32() * context.Persistance
}

// Dividing algorithm
func (g HeightMapGenerator) divide(context generateContext, stepSize int) {
	half := int(math.Floor(float64(stepSize) / 2.0))
	size := len(*context.HeightMap)
	wg := new(sync.WaitGroup)

	if half < 1 {
		return
	}

	for x := half; x < size; x += stepSize {
		for y := half; y < size; y += stepSize {
			wg.Add(1)
			go g.square(context, x, y, half, g.getOffset(stepSize, context), wg)
		}
	}

	wg.Wait()
	g.divide(context, half)
}

// "Square" part of algorithm
func (g HeightMapGenerator) square(context generateContext, x int, y int, size int, offset float32, wg *sync.WaitGroup) {
	defer wg.Done()
	a := g.getCellHeight(context, x-size, y-size, size)
	b := g.getCellHeight(context, x+size, y+size, size)
	c := g.getCellHeight(context, x-size, y+size, size)
	d := g.getCellHeight(context, x+size, y-size, size)

	average := (a + b + c + d) / 4
	(*context.HeightMap)[x][y] = average + offset

	g.diamond(context, x, y-size, size, g.getOffset(size, context))
	g.diamond(context, x-size, y, size, g.getOffset(size, context))
	g.diamond(context, x, y+size, size, g.getOffset(size, context))
	g.diamond(context, x+size, y, size, g.getOffset(size, context))
}

// "Diamond" part of alhorithm
func (g HeightMapGenerator) diamond(context generateContext, x int, y int, size int, offset float32) {
	a := g.getCellHeight(context, x, y-size, size)
	b := g.getCellHeight(context, x, y+size, size)
	c := g.getCellHeight(context, x-size, y, size)
	d := g.getCellHeight(context, x+size, y, size)

	average := (a + b + c + d) / 4

	(*context.HeightMap)[x][y] = average + offset
}

// Getting cell height. Random if outside height map boundaries
func (g HeightMapGenerator) getCellHeight(context generateContext, x int, y int, stepSize int) float32 {
	hm := *context.HeightMap
	if x >= len(hm) || x < 0 || y >= len(hm[x]) || y < 0 {
		return g.getOffset(stepSize, context)
	}
	return hm[x][y]
}

// Print for debug
func (heightMap HeightMap) Print() {
	for _, v := range heightMap {
		fmt.Printf("%v\n", v)
	}
}
