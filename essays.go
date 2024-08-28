package main

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/schollz/progressbar/v3"
)

type Worker struct {
	id      int
	Jobs    <-chan string
	Results chan<- string
}

type WorkerPool struct {
	Workers       []*Worker
	Jobs          chan string
	Results       chan string
	wg            sync.WaitGroup
	totalJobs     int
	jobCounter    int
	progressBar   *progressbar.ProgressBar
	progressMutex sync.Mutex
}

func NewWorkerPool(noWorkers int, noJobs int) *WorkerPool {
	// we create a channel that holds all the jobs that the workers need to complete
	jobs := make(chan string, noJobs)
	// after a worker has completed a job, a worker drops it off here
	results := make(chan string, noJobs)

	// this manages the workers, the queue and the result queue
	pool := &WorkerPool{
		// slice of workers that will be noWorkers long
		Workers:   make([]*Worker, noWorkers),
		Jobs:      jobs,
		Results:   results,
		totalJobs: noJobs,
	}
	// initialise our progress bar
	pool.progressBar = progressbar.Default(int64(noJobs))

	// create our workers
	for i := 0; i < noWorkers; i++ {
		worker := Worker{
			id: i,
			// all workers are looking at the same job queue and result queue
			Jobs:    jobs,
			Results: results,
		}
		// add our worker to the pool
		pool.Workers[i] = &worker
	}

	return pool
}

func (wp *WorkerPool) UpdateProgress() {
	wp.progressMutex.Lock()
	defer wp.progressMutex.Unlock()
	wp.jobCounter++
	wp.progressBar.Add(1)
}

func (wp *WorkerPool) Start() {
	for _, worker := range wp.Workers {
		wp.wg.Add(1)
		go worker.StartWork(&wp.wg, wp)
	}
}

func (wp *WorkerPool) Stop() {
	close(wp.Jobs)
	wp.wg.Wait()
	close(wp.Results)
}

func (wp *WorkerPool) AddJob(url string) {
	wp.Jobs <- url
}

// Worker method for processing Jobs
func (w *Worker) StartWork(wg *sync.WaitGroup, wp *WorkerPool) {
	defer wg.Done()

	// workers see the job queue, which is a list of URLS
	// once a worker has read from the channel, it processes and removes
	// Golang handles syncronisation for us regarding workers reading the same job
	for url := range w.Jobs {
		response, err := http.Get(url)
		if err != nil {
			w.Results <- fmt.Sprintf("Worker %d failed to fetch %s: %v", w.id, url, err)
			continue
		}

		body, err := io.ReadAll(response.Body)
		if body != nil {
			err = response.Body.Close()
		}

		if err != nil {
			w.Results <- fmt.Sprintf("Worker %d failed to read body of %s: %v", w.id, url, err)
			continue
		}

		// Simulate processing time
		time.Sleep(1 * time.Second)
		w.Results <- fmt.Sprintf("Counter Worker %d processed %s", w.id, url)
		wp.UpdateProgress()
	}
}

// processEssays processes a list of URLs using a worker pool.
func processEssays(urls []string) {
	const workers = 100
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
