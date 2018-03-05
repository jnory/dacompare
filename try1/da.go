package try1

import (
	"errors"

	"../common"
)

type NodeIndex int
type Branch int

type Slot struct {
	Base  NodeIndex
	Check NodeIndex
}

type DA []Slot

func isEmpty(da DA, n NodeIndex) bool {
	return da[n].Check < 0
}

func checkEmpty(da DA, n NodeIndex, c Branch) bool {
	next := da[n].Base + NodeIndex(c)
	return isEmpty(da, next)
}

func nextEmpty(da DA, n NodeIndex) NodeIndex {
	return NodeIndex(-da[n].Check)
}

func previousEmpty(da DA, n NodeIndex) NodeIndex {
	return -da[n].Base
}

func firstEmptySlot(da DA) NodeIndex {
	return -da[0].Check
}

func adjust(target NodeIndex, c Branch) NodeIndex {
	return target - NodeIndex(c)
}

func nextNodeWithoutCheck(da DA, n NodeIndex, c Branch) NodeIndex {
	return da[n].Base + NodeIndex(c)
}

func nextNode(da DA, n NodeIndex, c Branch) (NodeIndex, bool) {
	next := nextNodeWithoutCheck(da, n, c)
	if next < NodeIndex(len(da)) && da[next].Check == n {
		return next, true
	} else {
		return -1, false
	}
}

func realloc(da DA, n NodeIndex) DA {
	allocd := make(DA, 0, n)
	allocd = append(allocd, da...)
	for i := NodeIndex(len(da)); i < n; i++ {
		allocd[i].Base = NodeIndex(-(i-1))
		allocd[i].Check = NodeIndex(-(i+1))
	}
	return allocd
}

func getPrefixNodeIndex(da DA, prefix string) (NodeIndex, error) {
	i := NodeIndex(0)
	l := len(prefix)
	for j := 0; j < l; j++ {
		c := Branch(prefix[j])
		var success bool
		i, success = nextNode(da, i, c)
		if !success {
			return 0, errors.New("incorrect entry order")
		}
	}

	return i, nil
}

func getInitialBase(da DA, least Branch) NodeIndex {
	emptySlot := firstEmptySlot(da)
	candidate := adjust(emptySlot, least)
	for candidate < 0 {
		emptySlot = nextEmpty(da, emptySlot)
		candidate = adjust(emptySlot, least)
	}

	return candidate
}

func insert(da DA, prefix string, branches []Branch, n int) (DA, NodeIndex, error) {
	if n <= 0 {
		return da, 0, errors.New("no data")
	}

	prefixNode, err := getPrefixNodeIndex(da, prefix)
	if err != nil {
		return da, 0, err
	}

	da[prefixNode].Base = getInitialBase(da, branches[0])

	for {
		conflict := false
		for i := 0; i < n; i++ {
			if !checkEmpty(da, prefixNode, branches[i]) {
				conflict = true
				break
			}
		}

		if !conflict {
			break
		}
		currentHead := nextNodeWithoutCheck(da, prefixNode, branches[0])
		da[prefixNode].Base = adjust(nextEmpty(da, currentHead), branches[0])
	}

	largestNode := nextNodeWithoutCheck(da, prefixNode, branches[n - 1])
	if NodeIndex(len(da)) < largestNode {
		da = realloc(da, 2 * largestNode)
	}

	for i := 0; i < n; i++ {
		n := nextNodeWithoutCheck(da, prefixNode, branches[i])

		next := nextEmpty(da, n)
		previous := previousEmpty(da, n)

		da[previous].Check = -next
		da[next].Base = -previous

		da[n].Check = prefixNode
	}

	return da, largestNode, nil
}

func buildInitialList(da DA) {
	// head empty slot
	da[0].Check = -1

	for i := 1; i < len(da); i++ {
		// previous empty index
		da[i].Base = NodeIndex(-(i-1))

		// next empty index
		da[i].Check = NodeIndex(-(i+1))
	}
}

func NewDA(reader common.OrderedReader) (DA, error) {
	da := make(DA, reader.Size() * 2)
	buildInitialList(da)

	block := make([]Branch, 256)
	head := 0
	prefix := ""
	daSize := NodeIndex(0)
	for word := reader.Next(); !reader.EOD(); word = reader.Next() {
		l := len(word)
		if prefix != word[:l-1] && head != 0 {
			var err error
			var largestIndex NodeIndex
			da, largestIndex, err = insert(da, prefix, block, head)
			if err != nil {
				return nil, err
			}

			head = 0
			prefix = word[:l-1]

			if largestIndex > daSize {
				daSize = largestIndex
			}
		}

		block[head] = Branch(word[l-1])
		head++
	}

	return da[:daSize+1], nil
}