package main

import (
	"container/heap"
	"regexp"
	"strings"
	"sync"
)

type WordCounter struct {
	countMap map[string]int
	wordSet  map[string]struct{}
	heap     *minHeap
	mutex    sync.Mutex
}

// uses regex to check if the string only contains characters
func isAlphabetic(word string) bool {
	match, _ := regexp.MatchString("^[a-zA-Z]+$", word)
	return match
}

// counts the words in a given string if they match the criteria set
func (wc *WordCounter) CountWords(essay string) {
	wc.mutex.Lock()
	defer wc.mutex.Unlock()

	// Split the essay into words and loop over them
	for _, word := range strings.Fields(essay) {
		//lower case the word
		formatWord := strings.ToLower(word)

		// check if its all letters and more than 3 characters
		if len(formatWord) >= minWordLen && isAlphabetic(formatWord) {
			if _, ok := wc.wordSet[formatWord]; ok {
				wc.countMap[formatWord]++
			}
		}
	}
}

func (wc *WordCounter) GetTopKWords(k int) []WordCount {
	for word, count := range wc.countMap {
		heap.Push(wc.heap, WordCount{word, count})
		// If heap size exceeds k, remove the smallest item
		if wc.heap.Len() > k {
			heap.Pop(wc.heap)
		}

	}
	// Convert the heap to a slice and return it
	result := make([]WordCount, wc.heap.Len())
	for i := 0; wc.heap.Len() > 0; i++ {
		result[i] = heap.Pop(wc.heap).(WordCount)
	}
	return result
}

// seperate our word counter from our worker pools etc, it just initialises our heap, countMap and wordSet for our workers to use
func NewWordCounter(wordSet map[string]struct{}) *WordCounter {
	h := &minHeap{}
	heap.Init(h)
	return &WordCounter{
		heap:     h,
		countMap: make(map[string]int),
		wordSet:  wordSet,
	}
}
