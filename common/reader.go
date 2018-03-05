package common

import (
	"io/ioutil"
	"sort"
	"strings"
)

type OrderedReader interface {
	Next() string
	EOD() bool
	Size() int
}

func splitNodeAndBranch(word string) (node string, branch uint8) {
	return word[:len(word)-1], word[len(word)-1]
}


type orderedReader struct {
	data []string
	ptr int
}

func removeNoData(lines []string) []string{
	words := make([]string, 0, len(lines))

	for i := 0; i < len(lines); i++ {
		word := strings.Trim(lines[i], " ")
		if word != "" {
			words = append(words, word)
		}
	}

	return words
}

func compareNodeOrder(word1, word2 string) bool {
	node1, branch1 := splitNodeAndBranch(word1)
	node2, branch2 := splitNodeAndBranch(word2)

	if node1 == node2 {
		return branch1 < branch2
	}

	lenShorter := len(node1)
	if len(node2) < lenShorter {
		lenShorter = len(node2)
	}

	for k := 0; k < lenShorter; k++ {
		if node1[k] != node2[k] {
			return node1[k] < node2[k]
		}
	}

	return len(node1) < len(node2)
}

func detectShared(word1, word2 string) int {
	shorter := len(word1)
	if len(word2) < shorter {
		shorter = len(word2)
	}

	for i := 0; i < shorter; i++ {
		if word1[i] != word2[i] {
			return i
		}
	}

	return shorter
}

func insertNonExistingEntries(words []string) []string{
	sort.Strings(words)
	previous := ""
	data := make([]string, 0, len(words) * 2)
	for i := 0; i < len(words); i++ {
		current := words[i]

		for j := detectShared(previous, current) + 1; j < len(current); j++ {
			data = append(data, current[:j])
		}

		previous = current
		data = append(data, current)
	}

	return data
}

func NewOrderedReader(path string) (OrderedReader, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(data), "\n")
	words := removeNoData(lines)
	words = insertNonExistingEntries(words)

	sort.Slice(words, func(i, j int) bool {
		word1 := words[i]
		word2 := words[j]
		return compareNodeOrder(word1, word2)
	})

	return &orderedReader{
		data: words,
		ptr: 0,
	}, nil
}

func (r *orderedReader) Next() string {
	ret := r.data[r.ptr]
	r.ptr++
	return ret
}

func (r *orderedReader) EOD() bool {
	return len(r.data) <= r.ptr
}

func (r *orderedReader) Size() int {
	return len(r.data)
}
