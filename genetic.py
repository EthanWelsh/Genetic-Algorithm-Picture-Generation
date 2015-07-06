import time

from ga_drawing import *

seed_image = get_pic_from_url(
    'https://upload.wikimedia.org/wikipedia/commons/thumb/e/ea/Van_Gogh_-_Starry_Night_-_Google_Art_Project.jpg/1280px-Van_Gogh_-_Starry_Night_-_Google_Art_Project.jpg')


class Population:
    # What percentage of the next generation should new, made from crossover?
    crossover_rate = .7

    # What percentage of drawings should be mutated?
    mutation_rate = .02

    # Of the drawings that are mutated, how many shapes within said drawing should be altered?
    mutation_amount = .1

    def __init__(self, population_size, number_of_points_per_shape=3, max_shape_size=30, number_of_shapes=1000):
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
        return list(self._population.values())

    def update_spinner(self, list_size=10000):
        self.population_spinner = []
        population_fitness = sum([chromosome._fitness for chromosome in self.population()])

        for chromosome in self.population():
            fitness_in_population = chromosome._fitness / population_fitness
            spots_on_spinner = int(fitness_in_population * list_size)
            self.population_spinner.extend([chromosome.id for _ in range(0, spots_on_spinner)])

    def evolve(self, pct_child=0.7, pct_old=0.15, pct_new=0.15):
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
        min, avg, max = None, None, None
        population_fitness = 0

        for chromosome in self.population():
            min = chromosome._fitness if (min is None or min > chromosome._fitness) else min
            max = chromosome._fitness if (max is None or max < chromosome._fitness) else max
            population_fitness += chromosome._fitness

        avg = population_fitness / self.population_size
        return min, avg, max

    def best_chromosome(self):
        min_fitness = None
        best_drawing = None

        for chromosome in self.population():
            if (min_fitness is None or chromosome._fitness < min_fitness):
                min_fitness = chromosome._fitness
                best_drawing = chromosome.drawing

        return best_drawing.get_pic_rep()


class Chromosome:
    def __init__(self, id, width, height, number_of_points_per_shape, max_shape_size, number_of_shapes):
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
    pop = Population(population_size=100, number_of_points_per_shape=3, max_shape_size=10, number_of_shapes=1000)
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
