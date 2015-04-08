package main

import (
	"math"
	"math/rand"
	"image"
)

const (
	CROSSOVER_RATE  = .7
)

type Chromosome struct {
	pic  Pic
}

// Fitness function used to determine the degree of completion of the picture
func (chromosome *Chromosome) Score(original Pic) float64 {

	differenceInColors := 0.0

	worstScore := float64((math.MaxUint8 * 3) * (Width * Height))

	for r := 0; r < Width; r++ {
		for c := 0; c < Height; c++ {
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
		modifiedRow := randomInt(0, Width-1)
		modifiedCol := randomInt(0, Height-1)

		modifiedColor := randomInt(0, 2)

		red, green, blue := population[modifiedChromosome].pic.GetRGB(modifiedRow, modifiedCol)

		indexInByteToFlip := uint32(randomInt(0, 7))
		currentValueInByte := byte(0)

		if modifiedColor == 0 { // Red

			currentValueInByte = getBitFromByte(red, int(indexInByteToFlip))

			if currentValueInByte == 0 {
				red = setBitInByte(red, indexInByteToFlip, 1)
			} else {
				red = setBitInByte(red, indexInByteToFlip, 0)
			}


		} else if modifiedColor == 1 { // Green
			currentValueInByte = getBitFromByte(green, int(indexInByteToFlip))

			if currentValueInByte == 0 {
				green = setBitInByte(green, indexInByteToFlip, 1)
			} else {
				green = setBitInByte(green, indexInByteToFlip, 0)
			}
		} else { // Blue
			currentValueInByte = getBitFromByte(blue, int(indexInByteToFlip))

			if currentValueInByte == 0 {
				blue = setBitInByte(blue, indexInByteToFlip, 1)
			} else {
				blue = setBitInByte(blue, indexInByteToFlip, 0)
			}
		}

		population[modifiedChromosome].pic.SetRGB(modifiedRow, modifiedCol, red, green, blue)
	}

	return population
}

// Will perform a crossover operation between two chromosomes
func MateChromosome(a Chromosome, b Chromosome, chromosomeChan chan Chromosome) {

	var ra, ga, ba uint8
	var rb, gb, bb uint8
	var rc, gc, bc uint8
	var rd, gd, bd uint8
	var c, d Chromosome

	c.pic.img = *image.NewRGBA(image.Rect(0, 0, Width, Height))
	d.pic.img = *image.NewRGBA(image.Rect(0, 0, Width, Height))

	if rand.Float64() < CROSSOVER_RATE {

		for x := 0; x < Width; x++ {
			for y := 0; y < Height; y++ {

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
		chromosomeChan <- c
		chromosomeChan <- d
	}
	chromosomeChan <- a
	chromosomeChan <- b
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

	chromosome.pic.img = *image.NewRGBA(image.Rect(0, 0, Width, Height))


	for x := 0; x < Width; x++ {
		for y := 0; y < Height; y++ {

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
