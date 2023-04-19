package gocachebenchmarkplus_test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/Yiling-J/go-cache-benchmark-plus/clients"
)

var benchClients = []clients.Client[string, string]{
	&clients.Theine[string, string]{},
	&clients.Ristretto[string, string]{},
}

func BenchmarkGetParallel(b *testing.B) {
	keys := []string{}
	for i := 0; i < 100000; i++ {
		keys = append(keys, fmt.Sprintf("%d", i))
	}
	for _, client := range benchClients {
		client.Init(100000)
		b.ResetTimer()
		b.Run(client.Name(), func(b *testing.B) {
			b.RunParallel(func(p *testing.PB) {
				counter := 0
				for p.Next() {
					client.Get(keys[counter%100000])
					counter++
				}
			})
		})
		client.Close()
	}
}

func BenchmarkSetParallel(b *testing.B) {
	keys := []string{}
	for i := 0; i < 1000000; i++ {
		keys = append(keys, fmt.Sprintf("%d", i))
	}
	for _, client := range benchClients {
		client.Init(100000)
		b.ResetTimer()
		b.Run(client.Name(), func(b *testing.B) {
			b.RunParallel(func(p *testing.PB) {
				counter := 0
				for p.Next() {
					client.Set(keys[counter%1000000], "bar")
					counter++
				}
			})
		})
		client.Close()
	}

}

func BenchmarkZipfParallel(b *testing.B) {
	z := rand.NewZipf(rand.New(rand.NewSource(time.Now().UnixNano())), 1.0001, 10, 1000000)
	keys := []string{}
	for i := 0; i < 1000000; i++ {
		keys = append(keys, fmt.Sprintf("%d", z.Uint64()))
	}
	for _, client := range benchClients {
		client.Init(100000)
		b.ResetTimer()
		b.Run(client.Name(), func(b *testing.B) {
			b.RunParallel(func(p *testing.PB) {
				counter := 0
				for p.Next() {
					_, ok := client.Get(keys[counter%1000000])
					if !ok {
						client.Set(keys[counter%1000000], "bar")
					}
					counter++
				}
			})
		})
		client.Close()
	}
}
