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

// Internal struct for options
type generateOptions struct {
	Size        int
	Persistance float32
}

// Type of height map (2d array)
type HeightMap [][]float32

// Main method of HeightMapGenerator
// Generates HeightMap by options
func (g HeightMapGenerator) Generate(size int, persistance float32) HeightMap {
	options := generateOptions{
		Size:        int(math.Pow(2, float64(size))) + 1,
		Persistance: persistance,
	}
	heightMap := make(HeightMap, options.Size)
	for i := range heightMap {
		heightMap[i] = make([]float32, options.Size)
	}

	g.initCorners(&heightMap, options)
	g.divide(&heightMap, options, options.Size)

	return heightMap
}

// Initializing corners of height map
func (g HeightMapGenerator) initCorners(heightMap *HeightMap, opt generateOptions) {
	lastIndex := opt.Size - 1
	hm := *heightMap

	hm[0][0] = g.getOffset(opt.Size, opt)
	hm[0][lastIndex] = g.getOffset(opt.Size, opt)
	hm[lastIndex][0] = g.getOffset(opt.Size, opt)
	hm[lastIndex][lastIndex] = g.getOffset(opt.Size, opt)
}

// Getting random offset of height
func (g HeightMapGenerator) getOffset(stepSize int, opt generateOptions) float32 {
	return float32(stepSize) / float32(opt.Size) * rand.Float32() * opt.Persistance
}

// Dividing algorithm
func (g HeightMapGenerator) divide(heightMap *HeightMap, opt generateOptions, stepSize int) {
	half := int(math.Floor(float64(stepSize) / 2.0))
	size := len(*heightMap)
	wg := new(sync.WaitGroup)

	if half < 1 {
		return
	}

	for x := half; x < size; x += stepSize {
		for y := half; y < size; y += stepSize {
			wg.Add(1)
			go g.square(heightMap, opt, x, y, half, g.getOffset(stepSize, opt), wg)
		}
	}

	wg.Wait()
	g.divide(heightMap, opt, half)
}

// "Square" part of algorithm
func (g HeightMapGenerator) square(heightMap *HeightMap, opt generateOptions, x int, y int, size int, offset float32, wg *sync.WaitGroup) {
	defer wg.Done()
	a := g.getCellHeight(heightMap, opt, x-size, y-size, size)
	b := g.getCellHeight(heightMap, opt, x+size, y+size, size)
	c := g.getCellHeight(heightMap, opt, x-size, y+size, size)
	d := g.getCellHeight(heightMap, opt, x+size, y-size, size)

	average := (a + b + c + d) / 4
	(*heightMap)[x][y] = average + offset

	g.diamond(heightMap, opt, x, y-size, size, g.getOffset(size, opt))
	g.diamond(heightMap, opt, x-size, y, size, g.getOffset(size, opt))
	g.diamond(heightMap, opt, x, y+size, size, g.getOffset(size, opt))
	g.diamond(heightMap, opt, x+size, y, size, g.getOffset(size, opt))
}

// "Diamond" part of alhorithm
func (g HeightMapGenerator) diamond(heightMap *HeightMap, opt generateOptions, x int, y int, size int, offset float32) {
	a := g.getCellHeight(heightMap, opt, x, y-size, size)
	b := g.getCellHeight(heightMap, opt, x, y+size, size)
	c := g.getCellHeight(heightMap, opt, x-size, y, size)
	d := g.getCellHeight(heightMap, opt, x+size, y, size)

	average := (a + b + c + d) / 4

	(*heightMap)[x][y] = average + offset
}

// Getting cell height. Random if outside height map boundaries
func (g HeightMapGenerator) getCellHeight(heightMap *HeightMap, opt generateOptions, x int, y int, stepSize int) float32 {
	hm := *heightMap
	if x >= len(hm) || x < 0 || y >= len(hm[x]) || y < 0 {
		return g.getOffset(stepSize, opt)
	}
	return hm[x][y]
}

// Print for debug
func (heightMap HeightMap) Print() {
	for _, v := range heightMap {
		fmt.Printf("%v\n", v)
	}
}
