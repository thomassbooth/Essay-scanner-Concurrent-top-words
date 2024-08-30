package main

import (
	"container/heap"
	"reflect"
	"testing"
)

func TestIsAlphabetic(t *testing.T) {
	tests := []struct {
		word     string
		expected bool
	}{
		{"Hello", true},
		{"world", true},
		{"123", false},
		{"Hello123", false},
		{"hello-world", false},
		{"", false},
	}

	for _, test := range tests {
		if got := isAlphabetic(test.word); got != test.expected {
			t.Errorf("isAlphabetic(%q) = %v; want %v", test.word, got, test.expected)
		}
	}
}

func TestWordCounter_CountWords(t *testing.T) {
	wordSet := map[string]struct{}{
		"this":  {},
		"is":    {},
		"a":     {},
		"test":  {},
		"hello": {},
		"world": {},
	}

	wordCounter := NewWordCounter(wordSet)

	// Test case 1: basic counting
	essay := "This is a test Hello world"
	wordCounter.CountWords(essay)

	expectedCounts := map[string]int{
		"this":  1,
		"is":    0,
		"a":     0,
		"test":  1,
		"hello": 1,
		"world": 1,
	}

	for word, expectedCount := range expectedCounts {
		if count := wordCounter.countMap[word]; count != expectedCount {
			t.Errorf("unexpected count for word %q: got %d, want %d", word, count, expectedCount)
		}
	}

	// Test case 2: Case insensitivity and filtering non-alphabetic words
	essay2 := "HELLO World this is another TEST with numbers 123 and symbols #@!"
	wordCounter.CountWords(essay2)

	expectedCounts = map[string]int{
		"this":  2,
		"is":    0,
		"a":     0,
		"test":  2,
		"hello": 2,
		"world": 2,
	}

	for word, expectedCount := range expectedCounts {
		if count := wordCounter.countMap[word]; count != expectedCount {
			t.Errorf("unexpected count for word %q: got %d, want %d", word, count, expectedCount)
		}
	}
}

func TestWordCounter_GetTopKWords(t *testing.T) {
	wordSet := map[string]struct{}{
		"this":  {},
		"is":    {},
		"a":     {},
		"test":  {},
		"hello": {},
		"world": {},
	}

	wordCounter := NewWordCounter(wordSet)

	essay := "This is a test Hello world this is another test this is dave"
	wordCounter.CountWords(essay)

	// Test case: get top 3 words
	topWords := wordCounter.GetTopKWords(3)

	// is is not a valid word since its less than 3
	expectedTopWords := []WordCount{
		{"world", 1},
		{"test", 2},
		{"this", 3},
	}

	if !reflect.DeepEqual(topWords, expectedTopWords) {
		t.Errorf("unexpected top words: got %v, want %v", topWords, expectedTopWords)
	}

	// Test case: get top 5 words
	topWords = wordCounter.GetTopKWords(4)
	expectedTopWords = []WordCount{
		{"hello", 1},
		{"world", 1},
		{"test", 2},
		{"this", 3},
	}

	if !reflect.DeepEqual(topWords, expectedTopWords) {
		t.Errorf("unexpected top words: got %v, want %v", topWords, expectedTopWords)
	}
}

func TestNewWordCounter(t *testing.T) {
	wordSet := map[string]struct{}{
		"test":  {},
		"hello": {},
		"world": {},
	}

	wordCounter := NewWordCounter(wordSet)

	if wordCounter == nil {
		t.Fatal("NewWordCounter returned nil")
	}

	if wordCounter.countMap == nil {
		t.Error("countMap is nil")
	}

	if wordCounter.wordSet == nil {
		t.Error("wordSet is nil")
	}

	if wordCounter.heap == nil {
		t.Error("heap is nil")
	}

	// Test that the heap is properly initialized
	heap.Init(wordCounter.heap)
	if wordCounter.heap.Len() != 0 {
		t.Error("heap is not initialized to an empty heap")
	}
}
