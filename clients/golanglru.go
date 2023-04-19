package clients

import (
	"image/color"

	lru "github.com/hashicorp/golang-lru/v2"
	"gonum.org/v1/plot/vg/draw"
)

type LRU[K comparable, V any] struct {
	client *lru.Cache[K, V]
}

func (c *LRU[K, V]) Style() *Style {
	return &Style{
		Color: color.RGBA{R: 255, A: 255},
		Shape: draw.PyramidGlyph{},
	}
}

func (c *LRU[K, V]) Init(cap int) {
	client, err := lru.New[K, V](cap)
	if err != nil {
		panic(err)
	}
	c.client = client

}

func (c *LRU[K, V]) Get(key K) (V, bool) {
	return c.client.Get(key)
}

func (c *LRU[K, V]) Set(key K, value V) {
	c.client.Add(key, value)
}
func (c *LRU[K, V]) Name() string {
	return "lru"
}

func (c *LRU[K, V]) Close() {
}

type TwoQueue[K comparable, V any] struct {
	client *lru.TwoQueueCache[K, V]
}

func (c *TwoQueue[K, V]) Style() *Style {
	return &Style{
		Color: color.RGBA{R: 255, G: 178, B: 102, A: 255},
		Shape: draw.CrossGlyph{},
	}
}

func (c *TwoQueue[K, V]) Init(cap int) {
	client, err := lru.New2Q[K, V](cap)
	if err != nil {
		panic(err)
	}
	c.client = client

}

func (c *TwoQueue[K, V]) Get(key K) (V, bool) {
	return c.client.Get(key)
}

func (c *TwoQueue[K, V]) Set(key K, value V) {
	c.client.Add(key, value)
}
func (c *TwoQueue[K, V]) Name() string {
	return "2q"
}

func (c *TwoQueue[K, V]) Close() {
}

type Arc[K comparable, V any] struct {
	client *lru.ARCCache[K, V]
}

func (c *Arc[K, V]) Style() *Style {
	return &Style{
		Color: color.RGBA{R: 51, G: 51, B: 255, A: 255},
		Shape: draw.RingGlyph{},
	}
}

func (c *Arc[K, V]) Init(cap int) {
	client, err := lru.NewARC[K, V](cap)
	if err != nil {
		panic(err)
	}
	c.client = client

}

func (c *Arc[K, V]) Get(key K) (V, bool) {
	return c.client.Get(key)
}

func (c *Arc[K, V]) Set(key K, value V) {
	c.client.Add(key, value)
}
func (c *Arc[K, V]) Name() string {
	return "arc"
}

func (c *Arc[K, V]) Close() {
}
