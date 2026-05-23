package main

import (
	"fmt"
	"time"
	"sync"
	"math/rand"
)

func main() {
	fmt.Println("Running Ingestion Performance Test (Mock)...")

	const numEvents = 100000
	const numWorkers = 10

	start := time.Now()
	var wg sync.WaitGroup
	wg.Add(numWorkers)

	for i := 0; i < numWorkers; i++ {
		go func(workerID int) {
			defer wg.Done()
			for j := 0; j < numEvents/numWorkers; j++ {
				// Mock serialization and processing time
				_ = fmt.Sprintf(`{"event_id": "%d", "source": "perf-test", "type": "json", "payload": "..."}`, rand.Int())
			}
		}(i)
	}

	wg.Wait()
	duration := time.Since(start)

	fmt.Printf("Processed %d events in %v\n", numEvents, duration)
	fmt.Printf("Throughput: %.2f events/sec\n", float64(numEvents)/duration.Seconds())
}
