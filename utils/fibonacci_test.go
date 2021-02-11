//
// fibonacci_test.go
// Copyright (C) 2021 Toran Sahu <toran.sahu@yahoo.com>
//
// Distributed under terms of the MIT license.
//

package utils

import (
	"testing"
)

func TestFibonacciRange(t *testing.T) {
	fibOf5 := []int{0, 1, 1, 2, 3, 5}
	i := 0
	for f := range FibonacciRange(5) {
		if f != fibOf5[i] {
			t.Errorf("got %d, want %d", f, fibOf5[i])
		}
		i++
	}
}
