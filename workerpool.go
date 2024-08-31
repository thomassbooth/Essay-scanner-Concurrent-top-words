package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/cenkalti/backoff/v4"
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
	WordCounter   *WordCounter
	wg            sync.WaitGroup
	totalJobs     int
	jobCounter    int
	progressBar   *progressbar.ProgressBar
	progressMutex sync.Mutex
}

// create our worker pool
func NewWorkerPool(noWorkers int, noJobs int, wordCounter *WordCounter) *WorkerPool {
	// we create a channel that holds all the jobs that the workers need to complete
	jobs := make(chan string, noJobs)
	// after a worker has completed a job, a worker drops it off here
	results := make(chan string, noJobs)

	// this manages the workers, the queue and the result queue
	pool := &WorkerPool{
		// slice of workers that will be noWorkers long
		Workers:     make([]*Worker, noWorkers),
		Jobs:        jobs,
		Results:     results,
		totalJobs:   noJobs,
		progressBar: progressbar.Default(int64(noJobs)),
		WordCounter: wordCounter,
	}

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

// lock progress and add to the progress bar and jobCounter
func (wp *WorkerPool) UpdateProgress() {
	//lock our progress so we dont have multiple workers trying to update it at the same time
	wp.progressMutex.Lock()
	defer wp.progressMutex.Unlock()
	wp.jobCounter++
	wp.progressBar.Add(1)
}

// sping up all the workers on their different go routines
// making sure to add wait groups to ensure all workers are done before we close the results channel
func (wp *WorkerPool) Start() {
	for _, worker := range wp.Workers {
		wp.wg.Add(1)
		go worker.StartWork(&wp.wg, wp)
	}
}

// after adding all our jobs to the queue wait for them all to finish and close the channels
func (wp *WorkerPool) Stop() {
	close(wp.Jobs)
	wp.wg.Wait()
	close(wp.Results)
}

// add url to the job channel
func (wp *WorkerPool) AddJob(url string) {
	wp.Jobs <- url
}

// Worker method for processing Jobs
func (w *Worker) StartWork(wg *sync.WaitGroup, wp *WorkerPool) {
	defer wg.Done()

	for url := range w.Jobs {

		//we want an initial delay so were not hammering the external api and get locked out
		time.Sleep(backOffIntervals.FirstCallDelay)

		// request function to pass into the backoff.Retry function
		request := func() error {
			cleanedEssay, err, retry := fetchAndProcessEssay(url)
			if retry {
				return err // Retry the request if retry is true
			}
			//weve got our cleaned essay now we can count the words
			wp.WordCounter.CountWords(cleanedEssay)
			wp.UpdateProgress()
			return nil
		}

		// Create a new exponential backoff instance
		expBackoff := SetupBackoff(backOffIntervals)

		// Use the backoff.Retry function to handle retries
		err := backoff.Retry(request, backoff.WithMaxRetries(expBackoff, 5))
		if err != nil {
			w.Results <- fmt.Sprintf("Worker %d failed to process %s after retries: %v", w.id, url, err)
		} else {
			w.Results <- fmt.Sprintf("Worker %d processed %s successfully", w.id, url)
		}
	}
}

// processEssays processes a list of URLs using a worker pool.
func processEssays(urls []string) []WordCount {
	wordSet := GenerateWordSet(wordSetPath)
	wordCounter := NewWordCounter(wordSet)
	pool := NewWorkerPool(workers, len(urls), wordCounter)
	// Add jobs to the pool
	for _, url := range urls {
		pool.AddJob(url)
	}
	// Start all workers
	pool.Start()
	// Wait for all jobs to be processed
	pool.Stop()

	result := wordCounter.GetTopKWords(maxWords)

	return result
}
