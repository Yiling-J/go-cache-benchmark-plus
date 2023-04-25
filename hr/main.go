package main

import (
	"bufio"
	"compress/gzip"
	"encoding/binary"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Yiling-J/go-cache-benchmark-plus/clients"
	"golang.org/x/image/font"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

const (
	GET = "GET"
	SET = "SET"
)

var benchClients = []clients.Client[string, string]{
	&clients.Theine[string, string]{},
	&clients.Ristretto[string, string]{},
	&clients.LRU[string, string]{},
	&clients.TwoQueue[string, string]{},
	&clients.Arc[string, string]{},
}

type key struct {
	key string
	op  string
}

func zipfGen(keyChan chan key) {
	z := rand.NewZipf(rand.New(rand.NewSource(time.Now().UnixNano())), 1.0001, 10, 50000000)
	for i := 0; i < 1000000; i++ {
		keyChan <- key{key: fmt.Sprintf("key:%d", z.Uint64()), op: GET}
	}
	close(keyChan)
}

func ds1Gen(keyChan chan key) {
	f, err := os.Open("trace/ds1.trace.gz")
	if err != nil {
		panic(err)
	}
	gr, err := gzip.NewReader(f)
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(gr)
	for scanner.Scan() {
		s := strings.Split(scanner.Text(), " ")
		base, _ := strconv.Atoi(s[0])
		count, _ := strconv.Atoi(s[1])
		for i := 0; i < count; i++ {
			keyChan <- key{key: strconv.Itoa(base + i), op: GET}
		}
	}
	close(keyChan)
}

func s3Gen(keyChan chan key) {
	f, err := os.Open("trace/s3.trace.gz")
	if err != nil {
		panic(err)
	}
	gr, err := gzip.NewReader(f)
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(gr)
	for scanner.Scan() {
		s := strings.Split(scanner.Text(), " ")
		base, _ := strconv.Atoi(s[0])
		count, _ := strconv.Atoi(s[1])
		for i := 0; i < count; i++ {
			keyChan <- key{key: strconv.Itoa(base + i), op: GET}
		}
	}
	close(keyChan)
}

func scarabGen(keyChan chan key) {
	f, err := os.Open("trace/sc2.trace.gz")
	if err != nil {
		panic(err)
	}
	gr, err := gzip.NewReader(f)
	if err != nil {
		panic(err)
	}
	reader := bufio.NewReader(gr)
	for {
		buf := make([]byte, 8)
		_, err := io.ReadFull(reader, buf)
		if err != nil {
			close(keyChan)
			break
		}
		num := binary.BigEndian.Uint64(buf)
		keyChan <- key{key: strconv.Itoa(int(num)), op: GET}
	}

}

func fbGen(keyChan chan key) {
	f, err := os.Open("trace/fb.trace.gz")
	if err != nil {
		panic(err)
	}
	gr, err := gzip.NewReader(f)
	if err != nil {
		panic(err)
	}
	reader := csv.NewReader(gr)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			close(keyChan)
			break
		}
		if record[1] == "op" {
			continue
		}
		keyChan <- key{key: record[0], op: record[1]}
	}
}

func bench(client clients.Client[string, string], cap int, gen func(keyChan chan key)) float64 {
	counter := 0
	miss := 0
	done := false
	keyChan := make(chan key)
	go gen(keyChan)
	client.Init(cap)
	for !done {
		k, more := <-keyChan
		if more {
			counter++
			if counter%100000 == 0 {
				fmt.Print(".")
			}
			switch k.op {
			case GET:
				v, ok := client.Get(k.key)
				if ok {
					if v != k.key {
						panic("")
					}
				} else {
					miss++
					client.Set(k.key, k.key)
				}
			case SET:
				client.Set(k.key, k.key)
			}
		} else {
			done = true
		}
	}
	hr := float64(counter-miss) / float64(counter)
	fmt.Printf("\n--- %s %d hit ratio: %.3f\n", client.Name(), cap, hr)
	client.Close()
	time.Sleep(time.Second)
	return hr
}

type result struct {
	Cap   int
	Ratio float64
}

func newPlot(title string) *plot.Plot {
	p := plot.New()
	p.Title.Text = fmt.Sprintf("Hit Ratios - %s", title)
	p.X.Label.Text = "capacity"
	p.Y.Label.Text = "hit ratio"
	p.Legend.TextStyle.Font.Size = vg.Points(16)
	p.Legend.TextStyle.Font.Style = font.StyleOblique
	p.Title.TextStyle.Font.Size = vg.Points(16)
	p.Title.TextStyle.Font.Style = font.StyleOblique
	p.X.Label.TextStyle.Font.Size = vg.Points(14)
	p.Y.Label.TextStyle.Font.Size = vg.Points(14)
	return p

}

func updatePlot(plot *plot.Plot, client clients.Client[string, string], results []result) {
	data := plotter.XYs{}
	for _, r := range results {
		data = append(data, plotter.XY{X: float64(r.Cap), Y: r.Ratio})
	}
	line, points, err := plotter.NewLinePoints(data)
	if err != nil {
		panic(err)
	}
	style := client.Style()
	line.Color = style.Color
	points.Shape = style.Shape
	points.Radius = vg.Points(3)
	plot.Add(line, points)
	plot.Legend.Add(strings.ToLower(client.Name()), line, points)
}

func benchAndPlot(title string, caps []int, gen func(keyChan chan key)) {
	fmt.Printf("======= bench %s =======\n", strings.ToLower(title))
	p := newPlot(title)

	for _, client := range benchClients {
		var results []result
		cacheFile := fmt.Sprintf("results/%s-%s.data", client.Name(), title)
		cached, err := os.ReadFile(cacheFile)
		if err == nil {
			err = json.Unmarshal(cached, &results)
			if err != nil {
				panic(err)
			}
			fmt.Printf("cached result found: %s\n", cacheFile)
		} else {
			for _, cap := range caps {
				hr := bench(client, cap, gen)
				results = append(results, result{Cap: cap, Ratio: hr})
			}
			b, err := json.Marshal(results)
			if err != nil {
				panic(err)
			}
			err = os.WriteFile(cacheFile, b, 0644)
			if err != nil {
				panic(err)
			}
		}
		updatePlot(p, client, results)
	}

	if err := p.Save(
		16*vg.Inch, 9*vg.Inch, fmt.Sprintf("results/%s.png", strings.ToLower(title)),
	); err != nil {
		panic(err)
	}

}

func main() {

	benchAndPlot("Zipf", []int{100, 200, 500, 1000, 2000, 5000, 10000, 20000}, zipfGen)
	benchAndPlot("DS1", []int{1000000, 2000000, 3000000, 5000000, 6000000, 8000000}, ds1Gen)
	benchAndPlot("S3", []int{50000, 100000, 200000, 300000, 500000, 800000, 1000000}, s3Gen)
	benchAndPlot("SCARAB1H", []int{1000, 2000, 5000, 10000, 20000, 50000, 100000}, scarabGen)
	benchAndPlot("META", []int{10000, 20000, 50000, 80000, 100000}, fbGen)
}
