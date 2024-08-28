package main

import "fmt"

// processEssays processes a list of URLs using a worker pool.
func processEssays(urls []string) {
	pool := NewWorkerPool(workers, len(urls))

	// Start all workers
	pool.Start()

	// Add jobs to the pool
	for _, url := range urls {
		pool.AddJob(url)
	}

	// Wait for all jobs to be processed
	pool.Stop()

	// Process results
	for result := range pool.Results {
		fmt.Println(result)
	}
}
