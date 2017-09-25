package main

import (
	"encoding/xml"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"os"
)

const fntname = "killgo.fnt"

var killcolor = color.RGBA{0xf6, 0xe0, 0x53, 0xff}

func main() {
	fntFile, err := ioutil.ReadFile(fntname)
	if err != nil {
		panic("wtf")
	}
	fnt := struct {
		XMLName xml.Name `xml:"font"`
		Pages   struct {
			Page []struct {
				ID   int    `xml:"id,attr"`
				File string `xml:"file,attr"`
			} `xml:"page"`
		} `xml:"pages"`
		Chars struct {
			Char []struct {
				ID     int `xml:"id,attr"`
				X      int `xml:"x,attr"`
				Y      int `xml:"y,attr"`
				Height int `xml:"height,attr"`
				Width  int `xml:"width,attr"`
				Page   int `xml:"page,attr"`
			} `xml:"char"`
		} `xml:"chars"`
	}{}
	xml.Unmarshal(fntFile, &fnt)
	pages := make(map[int]image.Image)
	for _, page := range fnt.Pages.Page {
		f, err := os.Open(page.File)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		img, err := png.Decode(f)
		if err != nil {
			panic(err)
		}
		pages[page.ID] = img
	}
	for _, char := range fnt.Chars.Char {
		page := pages[char.Page]
		offsetX := (64 - char.Width) / 2
		dest := image.NewRGBA(image.Rect(0, 0, 64, 64))
		for x := 0; x < char.Width; x++ {
			for y := 0; y < 64; y++ {
				c := page.At(x+char.X, y+char.Y)
				r, g, b, a := c.RGBA()
				if a >= 0xff00 {
					r = multiple(r, killcolor.R)
					g = multiple(g, killcolor.G)
					b = multiple(b, killcolor.B)
				}
				dest.Set(x+offsetX, y, color.RGBA{byte(r >> 8), byte(g >> 8), byte(b >> 8), byte(a >> 8)})
			}
		}
		f, err := os.Create(fmt.Sprintf("../emoticons/klg%04x.png", char.ID))
		if err != nil {
			panic(err)
		}
		png.Encode(f, dest)
		defer f.Close()
	}
	fmt.Println("Done! Nice desu ne!")
}

func multiple(a uint32, b byte) uint32 {
	var c uint32
	for i := uint(0); i < 16; i++ {
		if (a>>i)%2 == 1 {
			c += uint32(b) << (i - 8)
		}
	}
	return c
}
