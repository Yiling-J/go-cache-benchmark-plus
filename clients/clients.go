package clients

import (
	"image/color"

	"gonum.org/v1/plot/vg/draw"
)

type Style struct {
	Color color.Color
	Shape draw.GlyphDrawer
}

type Client[K comparable, V any] interface {
	Init(cap int)
	Get(key K) (V, bool)
	Set(key K, value V)
	Name() string
	Style() *Style
	Close()
}
