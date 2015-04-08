package main

import (
	"math"
	"math/rand"
	//"fmt"
	//"image"
	"image"
	//"fmt"
	//"fmt"
)

const (
	CROSSOVER_RATE  = .7
)

type Chromosome struct {
	pic  Pic
	rows int
	cols int
}

// Fitness function used to determine the degree of completion of the picture
func (chromosome *Chromosome) Score(original Pic) float64 {

	differenceInColors := 0.0

	rows := original.img.Bounds().Dx()
	cols := original.img.Bounds().Dy()

	worstScore := float64((math.MaxUint8 * 3) * (rows * cols))

	for r := 0; r < original.img.Bounds().Dx(); r++ {
		for c := 0; c < original.img.Bounds().Dy(); c++ {
			ra, ga, ba := chromosome.pic.GetRGB(r, c)
			rorig, gorig, borig := original.GetRGB(r, c)

			differenceInColors += math.Abs(float64(ra-rorig)) + math.Abs(float64(ga-gorig)) + math.Abs(float64(ba-borig))
		}
	}

	return ((worstScore - differenceInColors)/worstScore)*100
}

// Will randomly mutate random genes in random chromosomes within a given population
func Mutate(population []Chromosome, chanceToModifyPopulation float64) []Chromosome {

	if chanceToModifyPopulation == 0 {
		return population
	}

	for rand.Float64() < chanceToModifyPopulation { // if you decided to mutate...

		modifiedChromosome := randomInt(0, len(population)) // pick a random chromosome to modify
		modifiedRow := randomInt(0, population[modifiedChromosome].rows-1)
		modifiedCol := randomInt(0, population[modifiedChromosome].cols-1)
		modifiedColor := randomInt(0, 2)

		mutantColorValue := uint8(randomInt(0, math.MaxUint8))

		red, green, blue := population[modifiedChromosome].pic.GetRGB(modifiedRow, modifiedCol)

		if modifiedColor == 0 { // Red
			red = mutantColorValue
		} else if modifiedColor == 1 { // Green
			green = mutantColorValue
		} else { // Blue
			blue = mutantColorValue
		}

		population[modifiedChromosome].pic.SetRGB(modifiedRow, modifiedCol, red, green, blue)
	}

	return population
}

// Will perform a crossover operation between two chromosomes
func MateChromosome(a Chromosome, b Chromosome) (c Chromosome, d Chromosome) {

	var ra, ga, ba uint8
	var rb, gb, bb uint8
	var rc, gc, bc uint8
	var rd, gd, bd uint8

	c.pic.img = *image.NewRGBA(image.Rect(0, 0, a.pic.img.Bounds().Dx(), a.pic.img.Bounds().Dy()))
	d.pic.img = *image.NewRGBA(image.Rect(0, 0, a.pic.img.Bounds().Dx(), a.pic.img.Bounds().Dy()))

	if rand.Float64() < CROSSOVER_RATE {

		for x := 0; x < a.rows; x++ {
			for y := 0; y < a.cols; y++ {

				pixelC := a.pic.img.RGBAAt(x, y)
				ra = pixelC.R
				ga = pixelC.G
				ba = pixelC.B

				pixelD := b.pic.img.RGBAAt(x, y)
				rb = pixelD.R
				gb = pixelD.G
				bb = pixelD.B

				rc, rd = crossBitString(ra, rb)
				gc, gd = crossBitString(ga, gb)
				bc, bd = crossBitString(ba, bb)

				pixelC.R = rc
				pixelC.G = gc
				pixelC.B = bc

				c.pic.img.Set(x, y, pixelC)

				pixelD.R = rd
				pixelD.G = gd
				pixelD.B = bd

				d.pic.img.Set(x, y, pixelD)
			}
		}
		return c, d
	}
	return a, b
}

func crossBitString(a, b uint8) (c, d uint8) {

	crossoverPoint := uint(randomInt(1, 7))

	var bitMask uint8 = 255 << crossoverPoint

	crossLeftc  :=  bitMask & a
	crossRightc := ^bitMask & b

	c = crossLeftc | crossRightc

	crossLeftd  := ^bitMask & a
	crossRightd :=  bitMask & b

	d = crossLeftd | crossRightd

	return c, d
}

// Generates a random gene sequence that represents a possible partial solution to the given board
func GetRandomChromosome(p *Pic) (chromosome Chromosome) {

	chromosome.rows = p.img.Bounds().Dx()
	chromosome.cols = p.img.Bounds().Dy()

	chromosome.pic.img = *image.NewRGBA(image.Rect(0, 0, p.img.Bounds().Dx(), p.img.Bounds().Dy()))


	for x := 0; x < chromosome.rows; x++ {
		for y := 0; y < chromosome.cols; y++ {

			randomRed := uint8(randomInt(0, math.MaxUint8))
			randomGreen := uint8(randomInt(0, math.MaxUint8))
			randomBlue := uint8(randomInt(0, math.MaxUint8))

			pixel := p.img.RGBAAt(x, y) // get the color at this pixel
			pixel.R = randomRed
			pixel.G = randomGreen
			pixel.B = randomBlue
			chromosome.pic.img.Set(x, y, pixel)
		}
	}

	return chromosome
}

// Generates a random integer between min and max (inclusive)
func randomInt(min int, max int) int {

	if min == max {
		return min
	}

	return rand.Intn(max) + min
}
