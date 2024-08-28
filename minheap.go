package main

import (
	"container/heap"
)

type WordCount struct {
	Word  string
	Count int
}

type WordHeap struct {
	pairs   []WordCount
	maxSize int
}

//implementation followed from https://pkg.go.dev/container/heap

func (h WordHeap) Len() int           { return len(h.pairs) }
func (h WordHeap) Less(i, j int) bool { return h.pairs[i].Count < h.pairs[j].Count }
func (h WordHeap) Swap(i, j int)      { h.pairs[i], h.pairs[j] = h.pairs[j], h.pairs[i] }

// Modifying the generic heap structure, we only push to the list if whats in there is less than 10
// or we kick out the smallest and replace with our new higher value
func (h *WordHeap) Push(x any) {
	//list isnt full yet
	if h.Len() < h.maxSize {
		h.pairs = append(h.pairs, x.(WordCount))
		//check if the count were trying to push on is greater than the smallest in the list
	} else if x.(WordCount).Count > h.pairs[0].Count {
		h.Pop() // Remove the smallest count
		h.pairs = append(h.pairs, x.(WordCount))
	}
	// need to reheapify the list
	heap.Fix(h, h.Len()-1)
}

// grab the top value, splice the rest from 1:n
func (h *WordHeap) Pop() any {
	old := h.pairs
	n := len(old)
	x := old[0]
	h.pairs = old[1:n]
	return x
}

// new heap instance
// h := &WordHeap{maxSize: 5}
// heap.Init(h)
