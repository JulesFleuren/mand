package main

import (
	"image/color"

	"github.com/boombuler/barcode"
)

// ColoredBarcode is a Barcode with a color, in order to be
// able to overwrite the At and ColorModel functions
type ColoredBarcode struct {
	barcode.Barcode
	mainColor       color.Color
	backgroundColor color.Color
}

func (c ColoredBarcode) At(x, y int) color.Color {
	if c.Barcode.At(x, y) == color.Black {
		return c.mainColor
	} else {
		return c.backgroundColor
	}
}

func (c ColoredBarcode) ColorModel() color.Model {
	return color.RGBAModel
}
