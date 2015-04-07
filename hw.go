package main

import "fmt"

func main() {

    img := Init("smiley.png")
    img.SetRGB(70,80,69,69,69)
    fmt.Println(img.GetRGB(70, 80))

}