import random
from io import BytesIO

from PIL import ImageStat
import requests
import PIL.ImageDraw
import PIL.Image
import PIL.ImageShow


class Drawing:
    class Shape:
        class Color:
            def __init__(self):
                self.r = random.randint(0, 255)
                self.g = random.randint(0, 255)
                self.b = random.randint(0, 255)
                self.a = int(255 * max(random.random() * random.random(), 0.2))

            def get_color_tup(self):
                """Return a tuple representing this color's RGBA values"""

                return tuple([self.r, self.g, self.b, self.a])

            def __str__(self):
                return "rgba({} {} {} {})".format(self.r, self.g, self.b, self.a)

        class Point:
            def __init__(self, restrictions):
                min_x, min_y, max_x, max_y = restrictions
                self.x = random.randint(min_x, max_x)
                self.y = random.randint(min_y, max_y)

            def __str__(self):
                return "({}, {})".format(self.x, self.y)

        def __init__(self, number_of_points, boundaries, max_size=None):
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
            if max_size is None:
                self.restrictions = (self.origin.x - min(abs(self.origin.x - min_x), max_size),  # min_x
                                     self.origin.y - min(abs(self.origin.y - min_y), max_size),  # min_y
                                     self.origin.x + min(abs(self.origin.x - max_x), max_size),  # max_x
                                     self.origin.y + min(abs(self.origin.y - max_y), max_size))  # max_y
            else:
                self.restrictions = (self.origin.x - abs(self.origin.x - min_x),  # min_x
                                     self.origin.y - abs(self.origin.y - min_y),  # min_y
                                     self.origin.x + abs(self.origin.x - max_x),  # max_x
                                     self.origin.y + abs(self.origin.y - max_y))  # max_y

            self.points = [self.Point(self.restrictions) for _ in range(0, self.number_of_points)]
            self.color = self.Color()

        def get_point_tup(self):
            """Returns a tuple of tuples of points in this shape"""

            return tuple(tuple([point.x, point.y]) for point in self.points)

        def replace_random_point(self):
            """Selects a random point within a given shape and replaces it with a new point"""

            index = random.randint(0, len(self.points) - 1)
            self.points[index] = self.Point(self.restrictions)

        def __str__(self):
            return "Color: {} Points: {}".format(self.color, [str(point) for point in self.points])

    def __init__(self, width, height, points_per_shape, max_shape_size = None, number_of_shapes=0):
        """
        :param width: Drawing width
        :param height: Drawing height
        :param points_per_shape: Number of points within each shape
        :param: max_shape_size: Largest amount of distance that points are allowed to fall from the origin
        :param: number_of_shapes: Number of shapes within each drawing
        """

        self.width = width
        self.height = height

        self.points_per_shape = points_per_shape
        self.max_shape_size = max_shape_size
        self.shapes = [Drawing.Shape(number_of_points=self.points_per_shape, max_size=self.max_shape_size,
                                     boundaries=(0, 0, self.width, self.height)) for _ in range(0, number_of_shapes)]

    def add_shape(self):
        """Add a random shape to the this drawing"""

        self.shapes.append(Drawing.Shape(self.points_per_shape, self.max_shape_size, (0, 0, self.width, self.height)))

    def get_pic_rep(self):
        """Returns the image representation version of this drawing"""

        img = PIL.Image.new("RGB", (self.width, self.height), color=(255, 255, 255))
        draw = PIL.ImageDraw.Draw(img, 'RGBA')

        for shape in self.shapes:
            draw.polygon(shape.get_point_tup(), fill=shape.color.get_color_tup())

        return img

    def show_image(self):
        """Display's this drawing as an image"""

        im = self.get_pic_rep()
        im.show()

    def closeness(self, img2):
        """"Measure's the closeness of this image to the seed image"""

        return sum(
            ImageStat.Stat(PIL.ImageChops.difference(self.get_pic_rep().convert("RGBA"), img2.convert("RGBA"))).sum)


def get_pic_from_url(url):
    """Requests the picture from the given URL and returns the Pillow representation of this image"""

    response = requests.get(url)
    return PIL.Image.open(BytesIO(response.content)).convert('RGBA')


def main():
    img1 = get_pic_from_url(
        'http://www.homedepot.com/catalog/productImages/400/f9/f901c92f-01d7-4bd4-b18e-464b705f92ad_400.jpg')
    img2 = get_pic_from_url('https://www.tonydiloreto.com/wp-content/uploads/2015/03/whitesquare.jpg')

    print(img1.size)
    print(img2.size)

    print(ImageStat.Stat(PIL.ImageChops.difference(img1.convert("RGBA"), img2.convert("RGBA"))).sum)


if __name__ == "__main__":
    main()
