package main

import (
	"math/rand"
)

const (
	CROSSOVER_RATE = .7
	CHROMOSOME_SIZE = 50
)

type Chromosome struct {
	genes [CHROMOSOME_SIZE]uint8
}

// Fitness function used to determine the degree of completion of the board
func (c Chromosome) Score() float64 {
	return 0.0
}

// Will randomly mutate random genes in random chromosomes within a given population
func Mutate(population []Chromosome, chanceToModifyPopulation float64) []Chromosome {

	if chanceToModifyPopulation == 0 {
		return population
	}

	for rand.Float64() < chanceToModifyPopulation { // if you decided to mutate...

		//modifiedChromosome := randomInt(0, len(population)) // pick a random chromosome to modify

	}

	return population
}

// Will perform a crossover operation between two chromosomes
func MateChromosome(a Chromosome, b Chromosome) (Chromosome, Chromosome) {

	if rand.Float64() < CROSSOVER_RATE {

	}
	return a, b
}

// Generates a random gene sequence that represents a possible partial solution to the given board
func GetRandomChromosome(p *Pic) (chromosome Chromosome) {

	return
}

// Returns the string representation of a particular chromosome
func (c *Chromosome) String() (ret string) {

	return ""
}

// Generates a random integer between min and max (inclusive)
func randomInt(min int, max int) int {

	if min == max {
		return min
	}

	return rand.Intn(max) + min
}