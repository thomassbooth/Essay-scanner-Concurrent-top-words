package main

import (
	"container/heap"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/cenkalti/backoff/v4"
	"github.com/schollz/progressbar/v3"
)

const workers = 2
const maxWords = 10

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
	topWords      *WordHeap
	topWordsMutex sync.Mutex
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

	h := &WordHeap{maxSize: maxWords}

	// this manages the workers, the queue and the result queue
	pool := &WorkerPool{
		// slice of workers that will be noWorkers long
		Workers:   make([]*Worker, noWorkers),
		Jobs:      jobs,
		topWords:  h,
		Results:   results,
		totalJobs: noJobs,
	}
	// initialise our progress bar
	pool.progressBar = progressbar.Default(int64(noJobs))
	heap.Init(pool.topWords)

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
	//lock our progress so we dont have multiple workers trying to update it at the same time
	wp.progressMutex.Lock()
	defer wp.progressMutex.Unlock()
	wp.jobCounter++
	wp.progressBar.Add(1)
}

func (wp *WorkerPool) PushWord(word string, count int) {
	wp.topWordsMutex.Lock()
	defer wp.topWordsMutex.Unlock()
	wp.topWords.Push(WordCount{Word: word, Count: count})
}

// sping up all the workers on their different go routines
// making sure to add wait groups to ensure all workers are done before we close the results channel
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

	for url := range w.Jobs {
		// Define a function to perform the HTTP request and parsing
		requestFunc := func() error {
			response, err := http.Get(url)
			if err != nil {
				return fmt.Errorf("failed to fetch %s: %v", url, err)
			}
			defer response.Body.Close()

			// Check the HTTP status code
			if response.StatusCode != http.StatusOK {
				fmt.Println(response.StatusCode, url)
				return fmt.Errorf("received non-200 response %d for %s", response.StatusCode, url)
			}

			body, err := io.ReadAll(response.Body)
			if err != nil {
				return fmt.Errorf("failed to read body of %s: %v", url, err)
			}

			text := ParseHTMLFile(string(body))
			wordCountMap := countWords(text)
			fmt.Println(wordCountMap)
			fmt.Println("hello")

			wp.UpdateProgress()
			return nil
		}

		// Create a new exponential backoff instance
		expBackoff := backoff.NewExponentialBackOff()

		// Use the backoff.Retry function to handle retries
		err := backoff.Retry(requestFunc, backoff.WithMaxRetries(expBackoff, 5))
		if err != nil {
			w.Results <- fmt.Sprintf("Worker %d failed to process %s after retries: %v", w.id, url, err)
		} else {
			w.Results <- fmt.Sprintf("Worker %d processed %s successfully", w.id, url)
		}
	}
}

// processEssays processes a list of URLs using a worker pool.
func processEssays(urls []string) {
	pool := NewWorkerPool(workers, len(urls))

	// Add jobs to the pool
	for _, url := range urls {
		pool.AddJob(url)
	}

	// Start all workers
	pool.Start()

	// Wait for all jobs to be processed
	pool.Stop()

	// Process results
	for result := range pool.Results {
		fmt.Println(result)
	}
}
