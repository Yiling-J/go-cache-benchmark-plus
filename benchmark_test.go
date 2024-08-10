package gocachebenchmarkplus_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/Yiling-J/go-cache-benchmark-plus/clients"
)

var benchClients = []clients.Client[string, string]{
	&clients.Theine[string, string]{},
	&clients.Ristretto[string, string]{},
	&clients.Otter[string, string]{},
}

func BenchmarkGetParallel(b *testing.B) {
	keys := []string{}
	for i := 0; i < 100000; i++ {
		keys = append(keys, fmt.Sprintf("%d", i))
	}
	total := len(keys) - 1

	for _, client := range benchClients {
		client.Init(total)
		for _, key := range keys {
			client.Set(key, key)
		}
		b.ResetTimer()
		b.Run(client.Name(), func(b *testing.B) {
			b.RunParallel(func(p *testing.PB) {
				counter := rand.Int() % total
				for p.Next() {
					client.Get(keys[counter%total])
					counter++
				}
			})
		})
		client.Close()
	}
}

func BenchmarkGetSingleParallel(b *testing.B) {
	keys := []string{}
	for i := 0; i < 100000; i++ {
		keys = append(keys, fmt.Sprintf("%d", i))
	}
	total := len(keys) - 1

	for _, client := range benchClients {
		client.Init(total)
		for _, key := range keys {
			client.Set(key, key)
		}
		b.ResetTimer()
		b.Run(client.Name(), func(b *testing.B) {
			b.RunParallel(func(p *testing.PB) {
				for p.Next() {
					client.Get(keys[0])
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
	total := len(keys) - 1
	for _, client := range benchClients {
		client.Init(100000)
		b.ResetTimer()
		b.Run(client.Name(), func(b *testing.B) {
			b.RunParallel(func(p *testing.PB) {
				counter := int(rand.Int() % total)
				for p.Next() {
					client.Set(keys[counter%total], "bar")
					counter++
				}
			})
		})
		client.Close()
	}

}
