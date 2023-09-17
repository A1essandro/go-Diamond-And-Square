package hm

import (
	"fmt"
	"math"
	"math/rand"
	"sync"
	"time"
)

func Hm() {
	println("asdasd")
}

type HeightMapGenerator struct {
}

type generateOptions struct {
	Size        int
	Persistance float32
}

type HeightMap [][]float32

func (g HeightMapGenerator) Generate(persistance float32, size int) HeightMap {
	options := generateOptions{
		Size:        int(math.Pow(2, float64(size))) + 1,
		Persistance: persistance,
	}
	heightMap := make(HeightMap, options.Size)
	for i := range heightMap {
		heightMap[i] = make([]float32, options.Size)
	}

	lastIndex := options.Size - 1
	heightMap[0][0] = g.getOffset(options.Size, options)
	heightMap[0][lastIndex] = g.getOffset(options.Size, options)
	heightMap[lastIndex][0] = g.getOffset(options.Size, options)
	heightMap[lastIndex][lastIndex] = g.getOffset(options.Size, options)

	g.divide(&heightMap, options, options.Size)

	return heightMap
}

func (g HeightMapGenerator) getOffset(stepSize int, opt generateOptions) float32 {
	return float32(stepSize) / float32(opt.Size) * rand.Float32() * opt.Persistance
}

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

func (g HeightMapGenerator) square(heightMap *HeightMap, opt generateOptions, x int, y int, size int, offset float32, wg *sync.WaitGroup) {
	defer wg.Done()
	time.Sleep(50 * time.Millisecond)
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

func (g HeightMapGenerator) diamond(heightMap *HeightMap, opt generateOptions, x int, y int, size int, offset float32) {
	a := g.getCellHeight(heightMap, opt, x, y-size, size)
	b := g.getCellHeight(heightMap, opt, x, y+size, size)
	c := g.getCellHeight(heightMap, opt, x-size, y, size)
	d := g.getCellHeight(heightMap, opt, x+size, y, size)

	average := (a + b + c + d) / 4

	(*heightMap)[x][y] = average + offset
}

func (g HeightMapGenerator) getCellHeight(heightMap *HeightMap, opt generateOptions, x int, y int, stepSize int) float32 {
	hm := *heightMap
	if x >= len(hm) || x < 0 || y >= len(hm[x]) || y < 0 {
		return g.getOffset(stepSize, opt)
	}
	return hm[x][y]
}

func (heightMap HeightMap) Print() {
	for _, v := range heightMap {
		fmt.Printf("%v\n", v)
	}
}
