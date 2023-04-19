package clients

import (
	"image/color"

	"github.com/Yiling-J/theine-go"
	"gonum.org/v1/plot/vg/draw"
)

type Theine[K comparable, V any] struct {
	client *theine.Cache[K, V]
}

func (c *Theine[K, V]) Style() *Style {
	return &Style{
		Color: color.RGBA{B: 255, A: 255},
		Shape: draw.BoxGlyph{},
	}
}

func (c *Theine[K, V]) Init(cap int) {
	client, err := theine.NewBuilder[K, V](int64(cap)).Build()
	if err != nil {
		panic(err)
	}
	c.client = client
}

func (c *Theine[K, V]) Get(key K) (V, bool) {
	return c.client.Get(key)
}

func (c *Theine[K, V]) Set(key K, value V) {
	c.client.Set(key, value, 1)
}
func (c *Theine[K, V]) Name() string {
	return "theine"
}

func (c *Theine[K, V]) Close() {
	c.client.Close()
}
