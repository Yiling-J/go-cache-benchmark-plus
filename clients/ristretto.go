package clients

import (
	"image/color"

	"github.com/dgraph-io/ristretto"
	"gonum.org/v1/plot/vg/draw"
)

type Ristretto[K comparable, V any] struct {
	client *ristretto.Cache
}

func (c *Ristretto[K, V]) Style() *Style {
	return &Style{
		Color: color.RGBA{G: 255, A: 255},
		Shape: draw.CircleGlyph{},
	}
}

func (c *Ristretto[K, V]) Init(cap int) {
	client, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: int64(cap * 10),
		MaxCost:     int64(cap),
		BufferItems: 64,
	})
	if err != nil {
		panic(err)
	}
	c.client = client
}

func (c *Ristretto[K, V]) Get(key K) (V, bool) {
	v, ok := c.client.Get(key)
	if ok {
		return v.(V), true
	}
	var zero V
	return zero, false
}

func (c *Ristretto[K, V]) Set(key K, value V) {
	c.client.Set(key, value, 1)
}
func (c *Ristretto[K, V]) Name() string {
	return "ristretto"
}

func (c *Ristretto[K, V]) Close() {
	c.client.Close()
}
