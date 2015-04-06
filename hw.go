package main

import "fmt"

func main() {

    img := Init("smiley.png")
    fmt.Println(img.GetRGB(70, 80))

}