package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"
)

const (
	CHANCE_TO_MUTATE_A_POPULATION = .90

	POPULATION_SIZE = 100
	UNASSIGNED      = 0
)

var original Pic

func main() {

	rand.Seed(int64(time.Now().UnixNano()))
	original = Init("smiley.png")
	population := getRandomPopulation()

	avg, max, min := getPopulationStats(population)
	fmt.Printf("AVG: %.5f MAX: %.5f MIN: %.5f\n", avg, max, min)


	for i := 0; i < 5; i++ {

		fname := fmt.Sprintf("results/res%d.png", i)
		EncodePNG(fname, &population[0].pic.img)

		population = evolve(population, 1, CHANCE_TO_MUTATE_A_POPULATION)

		avg, max, min = getPopulationStats(population)
		fmt.Printf("AVG: %.5f MAX: %.5f MIN: %.5f\n", avg, max, min)


	}



}

func getRandomPopulation() []Chromosome {
	population := make([]Chromosome, POPULATION_SIZE)

	fmt.Println("Generating", POPULATION_SIZE, "random solutions. This may take a while...")

	// Generate random partial solutions to the given board
	for i := range population {
		population[i] = GetRandomChromosome(&original)

		if i%(POPULATION_SIZE/10) == 0 {
			fmt.Print((float64(i)/POPULATION_SIZE)*100.00, "%    ")
		}
	}

	fmt.Println("100%\nDone generating solutions! Starting evolution...")
	return population
}

// Performs reproduction and mutations for a given number of iterations and returns the resulting population
func evolve(population []Chromosome, iterations int, chanceAtMutation float64) []Chromosome {

	for i := 0; i < iterations; i++ {

		population = getNextGeneration(population)
		//population = Mutate(population, chanceAtMutation)
	}

	fmt.Println(population[0].Score(original))

	return population

}

// Performs reproduction and returns the resulting population
func getNextGeneration(oldPopulation []Chromosome) (newPopulation []Chromosome) {

	var randomChromosomeSelector Spinner
	randomChromosomeSelector.addOptions(oldPopulation, original)

	newPopulation = make([]Chromosome, POPULATION_SIZE)

	for i := 0; i < len(newPopulation); i += 2 {

		// Get mating partner A & B
		phenotypeA := randomChromosomeSelector.Spin()
		phenotypeB := randomChromosomeSelector.Spin()

		// Mate them and add their children to the new population
		newPopulation[i], newPopulation[i+1] = MateChromosome(phenotypeA, phenotypeB)

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