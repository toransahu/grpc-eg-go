//
// stack.go
// Copyright (C) 2020 Toran Sahu <toran.sahu@yahoo.com>
//
// Distributed under terms of the MIT license.
//

package stack

type Stack []float32

func (s *Stack) IsEmpty() bool {
	return len(*s) == 0
}

func (s *Stack) Push(input float32) {
	*s = append(*s, input)
}

func (s *Stack) Pop() (float32, bool) {
	if s.IsEmpty() {
		return -1.0, false
	}
	item := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return item, true
}
