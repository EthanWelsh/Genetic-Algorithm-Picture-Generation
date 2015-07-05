import time
from ga_drawing import *

seed_image = get_pic_from_url(
    'https://upload.wikimedia.org/wikipedia/commons/thumb/e/ea/Van_Gogh_-_Starry_Night_-_Google_Art_Project.jpg/1280px-Van_Gogh_-_Starry_Night_-_Google_Art_Project.jpg')


class Population:
    def __init__(self, population_size, number_of_points_per_shape=3, max_shape_size=30, number_of_shapes=1000):

        # What percentage of the next generation should new, made from crossover?
        self.crossover_rate = .7

        # What percentage of drawings should be mutated?
        self.mutation_rate = .02

        # Of the drawings that are mutated, how many shapes within said drawing should be altered?
        self.mutation_amount = .1

        self.population_size = population_size
        width, height = seed_image.size

        self._population = {}
        for _ in range(0, self.population_size):
            chromosome_id = random.getrandbits(128)
            self._population[chromosome_id] = Chromosome(chromosome_id, width, height, number_of_points_per_shape,
                                                         max_shape_size, number_of_shapes)

        self.population_spinner = []

    def population(self):
        return list(self._population.values())

    def update_spinner(self, list_size=10000):

        self.population_spinner = []

        population_fitness = sum([chromosome.fitness for chromosome in self.population()])

        for chromosome in self.population():
            fitness_in_population = chromosome.fitness / population_fitness
            spots_on_spinner = int(fitness_in_population * list_size)
            self.population_spinner.extend([chromosome.id for _ in range(0, spots_on_spinner)])

    def evolve(self, pct_new_child=.8):

        new_population = []

        self.update_spinner()
        for _ in range(0, self.population_size):
            if random.uniform(0, 1) < pct_new_child:
                # make a new child and add it to the population
                a = self._population[random.choice(self.population_spinner)]
                b = self._population[random.choice(self.population_spinner)]
                new_population.append(Chromosome.mate(a, b))
            else:
                # pick an old organism to carry on to the next generation
                new_population.append(self._population[random.choice(self.population_spinner)])

        mutations_to_make = self.mutation_rate * self.population_size
        chromosomes_to_mutate = [random.choice(self.population()) for _ in range(0, int(mutations_to_make))]

        for chromosome in chromosomes_to_mutate:
            new_population.append(chromosome.mutate())

        self._population = {}

        for chromosome in new_population:
            chromosome_id = chromosome.id
            self._population[chromosome_id] = chromosome

        self.update_spinner()


class Chromosome:
    def __init__(self, id, width, height, number_of_points_per_shape, max_shape_size, number_of_shapes):
        self.drawing = Drawing(width, height, number_of_points_per_shape, max_shape_size, number_of_shapes)
        self.id = id
        self.fitness = self.fitness()

    def mutate(self):
        """Will randomly mutate random genes in random chromosomes within a given population"""

        num_of_shapes_to_mutate = int(Population.mutation_amount * len(self.drawing.shapes))

        for _ in range(0, num_of_shapes_to_mutate):
            if random.random() < .5:  # color mutation
                random.choice(self.drawing.shapes).color = Drawing.Shape.Color()
            else:  # points mutation
                pass


        # Recalculate organism fitness
        self.fitness = self.fitness()

        return self

    def fitness(self):
        return self.drawing.closeness(seed_image)

    @classmethod
    def mate(cls, chromosome_one, chromosome_two):
        """Returns a new chromosome which will be the child of the given two chromosome"""
        return chromosome_one


def main():
    start = time.clock()
    pop = Population(population_size=10, number_of_points_per_shape=3, max_shape_size=30, number_of_shapes=1000)
    print("Population Created: {}".format((time.clock() - start)))

    start = time.clock()
    pop.evolve()
    print("Population Evolved: {}".format((time.clock() - start)))

    start = time.clock()
    pop.evolve()
    print("Population Evolved: {}".format((time.clock() - start)))


if __name__ == '__main__':
    main()
