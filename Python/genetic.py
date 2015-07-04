class Population:
    pass


class Chromosome:
    CROSSOVER_RATE = .7

    def __init__(self):
        """Provides a totally random chromosome"""
        pass

    def mate(chromosome_one, chromosome_two):
        """Returns a new chromosome which will be the child of the given two chromosome"""
        pass

    def mutate(self):
        """Will randomly mutate random genes in random chromosomes within a given population"""