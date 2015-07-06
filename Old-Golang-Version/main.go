package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"
)

const (
	CHANCE_TO_MUTATE_A_POPULATION = .99

	POPULATION_SIZE = 1000
	UNASSIGNED      = 0
)

var original Pic
var Width  int
var Height int

func main() {

	rand.Seed(int64(time.Now().UnixNano()))
	original = Init("inputs/small.png")
	population := getRandomPopulation()

	for i := 0; i < 500; i++ {

		fname := fmt.Sprintf("results/res%d.png", i)

		avg, max, min := getPopulationStats(population)
		fmt.Printf("%d). AVG: %.5f MAX: %.5f MIN: %.5f\n", i, avg, max, min)

		maxScore := 0.0
		maxScoreIndex := 0

		population = evolve(population, 100, CHANCE_TO_MUTATE_A_POPULATION)

		for i := range population {

			if population[i].Score(original) > maxScore {
				maxScore = population[i].Score(original)
				maxScoreIndex = i
			}
		}

		EncodePNG(fname, &population[maxScoreIndex].pic.img)

	}
}

func getRandomPopulation() []Chromosome {
	population := make([]Chromosome, POPULATION_SIZE)

	fmt.Println("Generating", POPULATION_SIZE, "random solutions. This may take a while...")

	chromosomeChan := make(chan Chromosome, POPULATION_SIZE)

	// Generate random partial solutions to the given board
	go GetRandomChromosome(&original, chromosomeChan)

	for i := range population {
		population[i] = <-chromosomeChan
	}

	fmt.Println("100%\nDone generating solutions! Starting evolution...")
	return population
}

// Performs reproduction and mutations for a given number of iterations and returns the resulting population
func evolve(population []Chromosome, iterations int, chanceAtMutation float64) []Chromosome {

	for i := 0; i < iterations; i++ {

		population = getNextGeneration(population)
		population = Mutate(population, chanceAtMutation)
	}

	return population
}

// Performs reproduction and returns the resulting population
func getNextGeneration(oldPopulation []Chromosome) (newPopulation []Chromosome) {

	var randomChromosomeSelector Spinner
	randomChromosomeSelector.addOptions(oldPopulation, original)

	chromosomeChan := make(chan Chromosome, POPULATION_SIZE)
	spinChan       := make(chan Chromosome, POPULATION_SIZE)

	newPopulation = make([]Chromosome, POPULATION_SIZE)

	go randomChromosomeSelector.Spin(spinChan)

	for i := 0; i < POPULATION_SIZE/2; i++ {

		// Mate them and add their children to the new population
		go MateChromosome(<-spinChan, <-spinChan, chromosomeChan)

	}

	for i := 0; i < POPULATION_SIZE; i++ {
		newPopulation[i] = <- chromosomeChan
	}

	return newPopulation
}

// Provide the average, maximum, and minimum board scores in the population
func getPopulationStats(population []Chromosome) (avg float64, max float64, min float64) {

	var total float64 = 0
	var chromosomeScore float64 = 0

	max = 0
	min = math.MaxUint64

	for _, chromosome := range population {

		chromosomeScore = float64(chromosome.Score(original))

		total += chromosomeScore

		if chromosomeScore > max {
			max = chromosomeScore
		}

		if chromosomeScore < min {
			min = chromosomeScore
		}
	}

	avg = float64(total) / float64(len(population))

	return

}

// temporary timing function
func trace(s string) (string, time.Time) {
	log.Println("START:", s)

	return s, time.Now()
}

// temporary timing function
func un(s string, startTime time.Time) {
	endTime := time.Now()
	log.Println("  END:", s, "ElapsedTime in seconds:", endTime.Sub(startTime))
}
