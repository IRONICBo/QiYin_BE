package urltrie

import (
	"fmt"
	"reflect"
	"testing"
)

type testHook struct {
	Hook
	priority int64
	urls     []string
}

func (h *testHook) Priority() int64 {
	return h.priority
}

func (h *testHook) Patterns() []string {
	return h.urls
}

func TestUrlTrie(t *testing.T) {
	// test case
	testData := []struct {
		priority int64
		urls     []string
	}{
		{1, []string{"/gin"}},
		{2, []string{"/gin"}},
		{3, []string{"/gin/1"}},
	}

	trie := NewTrie()

	// range test data
	for _, data := range testData {
		trie.InsertBatch(data.urls, &testHook{
			urls:     data.urls,
			priority: data.priority,
		})
	}

	values, matched := trie.Match("/gin/1")
	if matched {
		fmt.Printf("Matched, values: %#v\n", reflect.ValueOf(values)) // Output: Matched, values: [1 2]
	} else {
		fmt.Println("No match found")
	}

	values, matched = trie.Match("/qiyin/v1")
	if matched {
		fmt.Printf("Matched, values: %#v\n", reflect.ValueOf(values)) // Output: Matched, values: [1 2]
	} else {
		fmt.Println("No match found")
	}
}
