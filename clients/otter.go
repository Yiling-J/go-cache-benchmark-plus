package clients

import (
	"image/color"

	"github.com/maypok86/otter"
	"gonum.org/v1/plot/vg/draw"
)

type Otter[K comparable, V any] struct {
	client *otter.Cache[K, V]
}

func (c *Otter[K, V]) Style() *Style {
	return &Style{
		Color: color.RGBA{G: 127, A: 127},
		Shape: draw.TriangleGlyph{},
	}
}

func (c *Otter[K, V]) Init(cap int) {
	client, err := otter.MustBuilder[K, V](cap).Build()
	if err != nil {
		panic(err)
	}
	c.client = &client
}

func (c *Otter[K, V]) Get(key K) (V, bool) {
	v, ok := c.client.Get(key)
	if ok {
		return v, true
	}
	var zero V
	return zero, false
}

func (c *Otter[K, V]) Set(key K, value V) {
	c.client.Set(key, value)
}
func (c *Otter[K, V]) Name() string {
	return "otter"
}

func (c *Otter[K, V]) Close() {
	c.client.Close()
}
