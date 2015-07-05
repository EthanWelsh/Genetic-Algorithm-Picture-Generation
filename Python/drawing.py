import random

import PIL.ImageDraw
import PIL.Image
import PIL.ImageShow


class Shape:
    class Color:
        def __init__(self):
            self.r = random.randint(0, 255)
            self.g = random.randint(0, 255)
            self.b = random.randint(0, 255)

        def get_color_tup(self):
            return tuple([self.r, self.g, self.b])

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

        # The boundaries in which this shapes' points must fall
        self.restrictions = (self.origin.x - min(abs(self.origin.x - min_x), max_size),  # min_x
                             self.origin.y - min(abs(self.origin.y - min_y), max_size),  # min_y
                             self.origin.x + min(abs(self.origin.x - max_x), max_size),  # max_x
                             self.origin.y + min(abs(self.origin.y - max_y), max_size))  # max_y

        self.points = [self.Point(self.restrictions) for _ in range(0, self.number_of_points)]
        self.color = self.Color()

    def __str__(self):
        return "Color: {} Points: {}".format(self.color, [str(point) for point in self.points])

    def get_point_tup(self):
        return tuple(tuple([point.x, point.y]) for point in self.points)

class Image:
    def __init__(self, width, height, points_per_shape, max_shape_size, number_of_shapes = 0):
        self.width = width
        self.height = height

        self.points_per_shape = points_per_shape
        self.max_shape_size = max_shape_size
        self.shapes = [Shape(self.points_per_shape, self.max_shape_size, (0, 0, self.width, self.height)) for _ in range(0, number_of_shapes)]

    def add_shape(self):
        self.shapes.append(Shape(self.points_per_shape, self.max_shape_size, (0, 0, self.width, self.height)))

    def get_pic_rep(self):
        im = PIL.Image.new("RGB", (self.width, self.height), color=(255, 255, 255))
        draw = PIL.ImageDraw.Draw(im)

        for shape in self.shapes:
            draw.polygon(shape.get_point_tup(), fill=shape.color.get_color_tup())
        return im

    def show_image(self):
        im = self.get_pic_rep()
        im.show()

    def closeness(self, other_image):
        return PIL.ImageChops.subtract(self, other_image)


img = Image(400, 400, 3, 15, 20)
img.show_image()
