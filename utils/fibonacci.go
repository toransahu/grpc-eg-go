//
// fibonacci.go
// Copyright (C) 2021 Toran Sahu <toran.sahu@yahoo.com>
//
// Distributed under terms of the MIT license.
//

package utils

func FibonacciRange(n int) <-chan int {
	ch := make(chan int)
	fn := make([]int, n+1, n+2)
	fn[0] = 0
	fn[1] = 1
	go func() {
		defer close(ch)
		for i := 0; i <= n; i++ {
			var f int
			if i < 2 {
				f = fn[i]
			} else {
				f = fn[i-1] + fn[i-2]
			}
			fn[i] = f
			ch <- f
		}
	}()
	return ch
}
