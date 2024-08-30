package main

type WordCount struct {
	Word  string
	Count int
}

type minHeap []WordCount

func (h minHeap) Len() int           { return len(h) }
func (h minHeap) Less(i, j int) bool { return (h)[i].Count < (h)[j].Count }
func (h minHeap) Swap(i, j int)      { (h)[i], (h)[j] = (h)[j], (h)[i] }

// just append to the list (go takes care of the shuffling)
func (h *minHeap) Push(x interface{}) {
	*h = append(*h, x.(WordCount))
}

// remove the value at the end of the array as this will be the smallest
func (h *minHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
