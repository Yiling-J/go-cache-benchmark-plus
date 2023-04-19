.PHONY: bench-hitratios bench-throughput

bench-ratios:
	go run ./hr

bench-throughput:
	go test -bench=. -run="^$$" -benchmem
