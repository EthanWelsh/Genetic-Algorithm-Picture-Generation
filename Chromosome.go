package main

import (
	"math"
	"math/rand"
)

const (
	CROSSOVER_RATE  = .7
	CHROMOSOME_SIZE = 50
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

	perfectScore := float64((math.MaxUint8 * 3) * (rows * cols))

	for r := 0; r < original.img.Bounds().Dx(); r++ {
		for c := 0; c < original.img.Bounds().Dy(); c++ {
			ra, ga, ba := chromosome.pic.GetRGB(r, c)
			rb, gb, bb := original.GetRGB(r, c)

			differenceInColors += math.Abs(float64(ra-rb)) + math.Abs(float64(ga-gb)) + math.Abs(float64(ba-bb))
		}
	}

	return perfectScore - differenceInColors
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

	var fromA uint8

	var ra, ga, ba uint8
	var rb, gb, bb uint8
	var rc, gc, bc uint8
	var rd, gd, bd uint8

	if rand.Float64() < CROSSOVER_RATE {

		for x := 0; x < a.rows; x++ {
			for y := 0; y < a.cols; y++ {

				ra, ga, ba = a.pic.GetRGB(x, y)
				rb, gb, bb = b.pic.GetRGB(x, y)

				fromA = uint8(randomInt(0, 1))
				if fromA == 1 {
					rc = ra
					rd = rb
				} else {
					rd = ra
					rc = rb
				}

				fromA = uint8(randomInt(0, 1))
				if fromA == 1 {
					gc = ga
					gd = gb
				} else {
					gd = ga
					gc = gb
				}

				fromA = uint8(randomInt(0, 1))
				if fromA == 1 {
					bc = ba
					bd = bb
				} else {
					bd = ba
					bc = bb
				}

				c.pic.SetRGB(x, y, rc, gc, bc)
				d.pic.SetRGB(x, y, rd, gd, bd)

			}
		}
		return c, d
	}
	return a, b
}

// Generates a random gene sequence that represents a possible partial solution to the given board
func GetRandomChromosome(p *Pic) (chromosome Chromosome) {

	chromosome.rows = p.img.Bounds().Dx()
	chromosome.cols = p.img.Bounds().Dy()

	for x := 0; x < chromosome.rows; x++ {
		for y := 0; y < chromosome.cols; y++ {

			randomRed := uint8(randomInt(0, math.MaxUint8))
			randomGreen := uint8(randomInt(0, math.MaxUint8))
			randomBlue := uint8(randomInt(0, math.MaxUint8))

			chromosome.pic.SetRGB(x, y, randomRed, randomGreen, randomBlue)
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
