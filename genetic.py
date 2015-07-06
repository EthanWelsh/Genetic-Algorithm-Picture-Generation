import time

from ga_drawing import *

seed_image = get_pic_from_url(
    'http://www.homedepot.com/catalog/productImages/400/f9/f901c92f-01d7-4bd4-b18e-464b705f92ad_400.jpg')


class Population:
    # What percentage of the next generation should new, made from crossover?
    crossover_rate = .7

    # What percentage of drawings should be mutated?
    mutation_rate = .02

    # Of the drawings that are mutated, how many shapes within said drawing should be altered?
    mutation_amount = .1

    def __init__(self, population_size, number_of_points_per_shape=3, max_shape_size=30, number_of_shapes=1000):
        """
        :param population_size: Number of drawings within this population
        :param number_of_points_per_shape: Number of points within each shape
        :param max_shape_size: Maximum amount of distance that each shape's points are allowed to fall outside of the origin
        :param number_of_shapes: Number of shapes within each drawing in the population
        """

        self.population_size = population_size
        self.width, self.height = seed_image.size

        self.number_of_points_per_shape = number_of_points_per_shape
        self.max_shape_size = max_shape_size
        self.number_of_shapes = number_of_shapes

        self._population = {}
        for _ in range(0, self.population_size):
            chromosome_id = random.getrandbits(128)
            self._population[chromosome_id] = Chromosome(chromosome_id, self.width, self.height,
                                                         self.number_of_points_per_shape, self.max_shape_size,
                                                         self.number_of_shapes)
        self.population_spinner = []

    def population(self):
        """Returns him a list of all the drawings within this population"""

        return list(self._population.values())

    def update_spinner(self, list_size=10000):
        """Fills the spinner for this population, measuring each chromosome's fitness and giving more spots to
        chromosomes with better fitness scores"""

        self.population_spinner = []
        population_fitness = sum([chromosome._fitness for chromosome in self.population()])

        for chromosome in self.population():
            fitness_in_population = chromosome._fitness / population_fitness
            spots_on_spinner = int(fitness_in_population * list_size)
            self.population_spinner.extend([chromosome.id for _ in range(0, spots_on_spinner)])

    def evolve(self, pct_child=0.6, pct_old=0.20, pct_new=0.20):
        """
        Evolve the population, applying crossover, mutation, and carry-over to form the new population
        :param pct_child: The percentage of chromosomes within the new population that are made from crossover.
        :param pct_old: Percentage of carry-over from the previous generation into the new generation.
        :param pct_new: Percentage of the new generation that is completely new and random.
        """

        new_population = []

        self.update_spinner()
        for _ in range(0, self.population_size):
            rand = random.uniform(0, pct_child + pct_old + pct_new)

            if rand < pct_child:
                # make a new child and add it to the population
                a = self._population[random.choice(self.population_spinner)]
                b = self._population[random.choice(self.population_spinner)]
                new_population.append(Chromosome.mate(a, b))
            elif rand < pct_child + pct_old:
                # pick an old organism to carry on to the next generation
                new_population.append(self._population[random.choice(self.population_spinner)])
            else:
                chromosome_id = random.getrandbits(128)
                new_population.append(
                    Chromosome(chromosome_id, self.width, self.height, self.number_of_points_per_shape,
                               self.max_shape_size, self.number_of_shapes))

        mutations_to_make = Population.mutation_rate * self.population_size
        chromosomes_to_mutate = [random.choice(self.population()) for _ in range(0, int(mutations_to_make))]

        for chromosome in chromosomes_to_mutate:
            new_population.append(chromosome.mutate())

        self._population = {}

        for chromosome in new_population:
            chromosome_id = chromosome.id
            self._population[chromosome_id] = chromosome

        self.update_spinner()

    def statistic(self):
        """Provides the minimum, average, and max fitness score for this population"""

        min, avg, max = None, None, None
        population_fitness = 0

        for chromosome in self.population():
            min = chromosome._fitness if (min is None or min > chromosome._fitness) else min
            max = chromosome._fitness if (max is None or max < chromosome._fitness) else max
            population_fitness += chromosome._fitness

        avg = population_fitness / self.population_size
        return min, avg, max

    def best_chromosome(self):
        """Finds the drawing that is most representative of the seed image and returns that drawing as an image"""

        min_fitness = None
        best_drawing = None

        for chromosome in self.population():
            if (min_fitness is None or chromosome._fitness < min_fitness):
                min_fitness = chromosome._fitness
                best_drawing = chromosome.drawing

        return best_drawing.get_pic_rep()


class Chromosome:
    def __init__(self, id, width, height, number_of_points_per_shape, max_shape_size, number_of_shapes):
        """
        :param id: This chromosome's unique identifier id
        :param width: The drawing's width within this chromosome
        :param height: The drawing's height within this chromosome
        :param number_of_points_per_shape: Number of points within each shape in the drawing within this chromosome
        :param max_shape_size: The maximum amount of distance that points are allowed to fall away from a shape's origin
        :param number_of_shapes: The number of shapes within each drawing
        """

        self.drawing = Drawing(width, height, number_of_points_per_shape, max_shape_size, number_of_shapes)
        self.id = id
        self._fitness = self.fitness() if number_of_shapes > 0 else 0

    def mutate(self):
        """Will randomly mutate random genes in random chromosomes within a given population"""

        num_of_shapes_to_mutate = int(Population.mutation_amount * len(self.drawing.shapes))

        for _ in range(0, num_of_shapes_to_mutate):
            mutate_me = random.choice(self.drawing.shapes)

            if random.random() < .5:  # color mutation
                mutate_me.color = Drawing.Shape.Color()
            else:  # points mutation
                mutate_me.replace_random_point()

        # Recalculate organism fitness
        self._fitness = self.fitness()

        return self

    def fitness(self):
        """Return the difference between the seed_image and this chromosome's drawing. Lower scores are better."""

        return self.drawing.closeness(seed_image)

    @classmethod
    def mate(cls, chromosome_one, chromosome_two):
        """Returns a new chromosome which will be the child of the given two chromosome"""

        new_chromosome = Chromosome(random.getrandbits(128), chromosome_one.drawing.width,
                                    chromosome_one.drawing.height, 0, 0, 0)

        for i in range(0, len(chromosome_one.drawing.shapes)):
            if random.random() < .5:
                new_chromosome.drawing.shapes.append(chromosome_one.drawing.shapes[i])
            else:
                new_chromosome.drawing.shapes.append(chromosome_two.drawing.shapes[i])

        new_chromosome._fitness = new_chromosome.fitness()
        return new_chromosome


def main():
    start = time.clock()
    pop = Population(population_size=100, number_of_points_per_shape=3, max_shape_size=5, number_of_shapes=500)
    print("Population Created: {}".format((time.clock() - start)))

    print(pop.statistic())
    pop.best_chromosome().save("out/before.png", "PNG")

    for i in range(0, 100):
        if i % 10 == 0:
            print(pop.statistic())
            pop.best_chromosome().save("out/{}.png".format(i / 10), "PNG")
        pop.evolve()

    print(pop.statistic())
    pop.best_chromosome().save("out/after.png", "PNG")


if __name__ == '__main__':
    main()
