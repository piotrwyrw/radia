package vtheme

import (
	"image/color"

	"fyne.io/fyne/v2"
	ftheme "fyne.io/fyne/v2/theme"
)

type RadiaTheme struct {
	Fallback fyne.Theme
}

func rgb(r, g, b uint8) color.Color {
	return color.RGBA{R: r, G: g, B: b, A: 255}
}

func (rt RadiaTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case ftheme.ColorNameBackground:
		return rgb(10, 10, 10)
	case ftheme.ColorNameForeground:
		return rgb(255, 250, 251)
	case ftheme.ColorNameError:
		return rgb(231, 76, 60)
	case ftheme.ColorNamePrimary:
		return rgb(0, 200, 255)
	case ftheme.ColorNameShadow:
		return rgb(0, 0, 0)
	case ftheme.ColorNameInputBackground:
		return rgb(15, 15, 15)
	case ftheme.ColorNameInputBorder:
		return rgb(50, 50, 50)
	default:
		return rt.Fallback.Color(name, variant)
	}
}

func (rt RadiaTheme) Font(style fyne.TextStyle) fyne.Resource {
	return rt.Fallback.Font(style)
}

func (rt RadiaTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return rt.Fallback.Icon(name)
}

func (rt RadiaTheme) Size(name fyne.ThemeSizeName) float32 {
	return rt.Fallback.Size(name)
}
