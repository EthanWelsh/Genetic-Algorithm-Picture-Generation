package main

import (
	"math"
	"math/rand"
	"image"
	"fmt"
)

const (
	CROSSOVER_RATE  = .7
)

type Chromosome struct {
	pic  Pic
}

type MyColor struct {
	red uint8
	green uint8
	blue uint8
}



// Fitness function used to determine the degree of completion of the picture
func (chromosome *Chromosome) Score(original Pic) float64 {

	var worstScore float64 = float64(distance(0.0, 0.0, 0.0, 255.0, 255.0, 255.0) * float64(Width * Height))

	differenceInColors := 0.0

	for r := 0; r < Width; r++ {
		for c := 0; c < Height; c++ {
			ra, ga, ba := chromosome.pic.GetRGB(r, c)
			rorig, gorig, borig := original.GetRGB(r, c)
			differenceInColors += distance(ra, ga, ba, rorig, gorig, borig)
		}
	}

	return ((worstScore - differenceInColors)/worstScore)*100
}

func distance(x1, y1, z1, x2, y2, z2 uint8) float64 {
	return math.Sqrt(math.Pow(float64(x2-x1), 2) + math.Pow(float64(y2-y1), 2) + math.Pow(float64(z2-z1), 2))
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

	var aa MyColor
	var bb MyColor

	var chromC, chromD Chromosome

	chromC.pic.img = *image.NewRGBA(image.Rect(0, 0, Width, Height))
	chromD.pic.img = *image.NewRGBA(image.Rect(0, 0, Width, Height))

	var rand uint8

	for r := 0; r < Width; r++ {
		for c := 0; c < Height; c++ {

			aa.red, aa.green, aa.blue = a.pic.GetRGB(r, c)
			bb.red, bb.green, bb.blue = b.pic.GetRGB(r, c)

			rand = uint8(randomInt(0, 24))

			pixelC, pixelD := crossPixel(rand, aa, bb)

			chromC.pic.SetRGB(r, c, pixelC.red, pixelC.green, pixelC.blue)
			chromD.pic.SetRGB(r, c, pixelD.red, pixelD.green, pixelD.blue)


		}
	}

	chromosomeChan <- chromC
	chromosomeChan <- chromD

}

func crossPixel(split uint8, a MyColor, b MyColor) (c MyColor, d MyColor) {

	if split < 0 || split > 24 {
		return
	} else if split < 8 {
    	c.red, d.red = cut(split, a.red, b.red)
		c.green, d.green = swap(a.green, b.green)
		c.blue, d.blue = swap(a.blue, b.blue)
	} else if split < 16 {
		c.red, d.red = swap(a.red, b.red)
		c.green, d.green = cut(split - 8, a.green, b.green)
		c.blue, d.blue = swap(a.blue, b.blue)
	} else if split < 24 {
		c.red, d.red = swap(a.red, b.red)
		c.blue, d.blue = swap(a.blue, b.blue)
		c.green, d.green = cut(split - 24, a.green, b.green)
	}

	return
}

func swap(a uint8, b uint8) (uint8, uint8) {
	return b, a
}

func cut(split uint8, a uint8, b uint8) (c uint8, d uint8) {
	amask := uint8(0xff >> split)
	bmask := uint8(^amask)

	c = (a & amask) | (b & bmask)
	d = (b & amask) | (a & bmask)

	return
}

// Generates a random gene sequence that represents a possible partial solution to the given board
func GetRandomChromosome(p *Pic, chromosomeChan chan Chromosome) {

	var chromosome Chromosome

	for i := 0; i < POPULATION_SIZE; i++ {
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

		chromosomeChan <- chromosome

		if i%(POPULATION_SIZE/10) == 0 {
			fmt.Print((float64(i)/POPULATION_SIZE)*100.00, "%    ")
		}
	}
}

// Generates a random integer between min and max (inclusive)
func randomInt(min int, max int) int {

	if min == max {
		return min
	}

	return rand.Intn(max) + min
}
