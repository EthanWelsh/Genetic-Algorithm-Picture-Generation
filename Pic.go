package main

import (
    "bufio"
    //"flag"
    //"fmt"
    "image"
    //"image/color"
    "image/draw"
    "image/png"
    //"io/ioutil"
    "log"
    "os"
)

type Pic struct {
    img image.RGBA
}


func Init(pictureInputFile string) Pic {
    return Pic(imageToRGBA(decodeImage(pictureInputFile)))
}

func (p Pic) GetRGB(x, y int) int {

    return p.img.RGBAAt(x, y) // get the color at this pixel
}

// will write out a given image to a given path in filename
func encodePNG(filename string, img image.Image) {
    fo, err := os.Create(filename)

    if err != nil {
        log.Fatalf("Error creating file %s: %v", filename, err)
    }

    defer fo.Close()
    defer fo.Sync()

    writer := bufio.NewWriter(fo)
    defer writer.Flush()

    err = png.Encode(writer, img)
}

// convert given image to RGBA image
func imageToRGBA(src image.Image) *image.RGBA {
    b := src.Bounds()

    var m *image.RGBA
    var width, height int

    width = b.Dx()
    height = b.Dy()

    m = image.NewRGBA(image.Rect(0, 0, width, height))
    draw.Draw(m, m.Bounds(), src, b.Min, draw.Src)
    return m
}

// read and return an image at the given path
func decodeImage(filename string) image.Image {
    inFile, err := os.Open(filename)

    if err != nil {
        log.Fatalf("Error opening file %s: %v", filename, err)
    }

    defer inFile.Close()
    reader := bufio.NewReader(inFile)
    img, _, err := image.Decode(reader)
    return img
}

// given a bit will return a bit from that byte
func getBitFromByte(b byte, indexInByte int) byte {
    b = b << uint(indexInByte)
    var mask byte = 0x80

    var bit byte = mask & b

    if bit == 128 {
        return 1
    }
    return 0
}

// sets a specific bit in a byte to a given value and returns the new byte
func setBitInByte(b byte, indexInByte uint32, bit byte) byte {
    var mask byte = 0x80
    mask = mask >> uint(indexInByte)

    if bit == 0 {
        mask = ^mask
        b = b & mask
    } else if bit == 1 {
        b = b | mask
    }
    return b
}
