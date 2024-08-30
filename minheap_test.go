package main

import (
	"container/heap"
	"testing"
)

func TestMinHeap_PushAndPop(t *testing.T) {
	h := &minHeap{}
	heap.Init(h)

	// Push items into the heap
	heap.Push(h, WordCount{"a", 5})
	heap.Push(h, WordCount{"b", 3})
	heap.Push(h, WordCount{"c", 8})
	heap.Push(h, WordCount{"d", 1})

	if h.Len() != 4 {
		t.Errorf("Expected heap length to be 4, got %d", h.Len())
	}

	// Pop the smallest item
	smallest := heap.Pop(h).(WordCount)
	if smallest.Count != 1 {
		t.Errorf("Expected smallest count to be 1, got %d", smallest.Count)
	}

	// Verify the new smallest item
	newSmallest := heap.Pop(h).(WordCount)
	if newSmallest.Count != 3 {
		t.Errorf("Expected next smallest count to be 3, got %d", newSmallest.Count)
	}

	// Check heap length after pops
	if h.Len() != 2 {
		t.Errorf("Expected heap length to be 2 after pops, got %d", h.Len())
	}
}

func TestMinHeap_EmptyPop(t *testing.T) {
	h := &minHeap{}
	heap.Init(h)

	// Ensure popping from an empty heap behaves correctly
	if len(*h) != 0 {
		t.Errorf("Expected heap to be empty, got length %d", len(*h))
	}

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic when popping from an empty heap, but did not panic")
		}
	}()

	heap.Pop(h)
}
