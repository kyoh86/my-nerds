package model

import (
	"image"
	"image/color"
)

type Page struct {
	Rel    string
	Format string
	Config image.Config
}

func ColorModelToString(cm color.Model) string {
	switch cm {
	case color.RGBAModel:
		return "RGBA"
	case color.RGBA64Model:
		return "RGBA64"
	case color.NRGBAModel:
		return "NRGBA"
	case color.NRGBA64Model:
		return "NRGBA64"
	case color.AlphaModel:
		return "Alpha"
	case color.Alpha16Model:
		return "Alpha16"
	case color.GrayModel:
		return "Gray"
	case color.Gray16Model:
		return "Gray16"
	case color.YCbCrModel:
		return "YCbCr"
	case color.NYCbCrAModel:
		return "NYCbCrA"
	case color.CMYKModel:
		return "CMYK"
	}
	if _, ok := cm.(color.Palette); ok {
		return "Palette"
	}
	return "Unknown"
}
