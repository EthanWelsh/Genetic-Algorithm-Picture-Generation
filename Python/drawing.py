import random
import PIL.ImageDraw
import PIL.Image
import PIL.ImageShow
import sys


class Shape:
    class Color:
        def __init__(self):
            self.r = random.randint(0, 255)
            self.g = random.randint(0, 255)
            self.b = random.randint(0, 255)

        def distance(self, other_color):
            return abs(self.r - other_color.r) + \
                   abs(self.g - other_color.g) + \
                   abs(self.b - other_color.b)

        def __str__(self):
            return "{} {} {}".format(self.r, self.g, self.b)

    class Point:
        def __init__(self, restrictions):
            min_x, min_y, max_x, max_y = restrictions
            self.x = random.randint(min_x, max_x)
            self.y = random.randint(min_y, max_y)

        def __str__(self):
            return "({}, {})".format(self.x, self.y)

    def __init__(self, number_of_points, max_size, boundaries):
        """
        :param number_of_points: Integer or tuple of (min, max) number of points allowed in shape. Inclusive.
        :param max_size: The maximum allowed size of the shape, measured in a box around the shape.
        :param boundaries: Tuple of (min_x, min_y, max_x, max_y) that describes the boundaries that this shape must
        fall inside.
        """

        if type(number_of_points) is tuple:
            min_num_of_points, max_num_of_points = number_of_points
            self.number_of_points = random.randint(min_num_of_points, max_num_of_points)
        elif type(number_of_points) is int:
            self.number_of_points = number_of_points
        else:
            raise TypeError

        self.origin = self.Point(boundaries)
        min_x, min_y, max_x, max_y = boundaries

        # The relative points inside of points must fall within these restrictions
        self.restrictions = (min(abs(self.origin.x - min_x), abs(self.origin.x - max_size)),
                             min(abs(self.origin.y - min_y), abs(self.origin.y - max_size)),
                             min(abs(self.origin.x - max_x), abs(self.origin.x + max_size)),
                             min(abs(self.origin.y - max_y), abs(self.origin.y + max_size)))

        self.points = [self.Point(self.restrictions) for _ in range(0, self.number_of_points)]
        self.color = self.Color()

    def __str__(self):
        return "Color: {} Points: {}".format(self.color, self.points)


class Image:
    def __init__(self, width, height):
        self.points = []
        self.width = width
        self.height = height

    def add_point(self, point):
        self.points.append(point)

    def get_pic_rep(self):
        im = PIL.Image.new("RGB", (self.width, self.height), color=(255, 255, 255))
        draw = PIL.ImageDraw.Draw(im)
        draw.polygon(self.points, fill=0)
        return im

    def show_image(self):
        im = self.get_pic_rep()
        im.show()

    def closeness(self, other_image):
        return PIL.ImageChops.subtract(self, other_image)

number_of_points = 3
max_size = 10
boundaries = (0, 0, 20, 20)

s = Shape(number_of_points, max_size, boundaries)


#img = Image(400, 400)
#img.add_point((50, 100))
#img.add_point((300, 700))
#img.add_point((15, 40))
#img.show_image()
