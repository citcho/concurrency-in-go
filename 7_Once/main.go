package main

import (
	"fmt"
	"sync"
)

func main() {
	var count int
	increment := func() {
		count++
	}
	var once sync.Once
	var increments sync.WaitGroup

	increments.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			defer increments.Done()
			once.Do(increment)
		}()
	}

	increments.Wait()
	fmt.Printf("count is %d\n", count) // count is 1

	// ------------------------------------------------------

	var count2 int
	increment2 := func() {
		count2++
	}
	decrement2 := func() {
		count2--
	}

	var once2 sync.Once
	once2.Do(increment2)
	once2.Do(decrement2)

	fmt.Printf("count is %d\n", count2) // count is 1

	// ------------------------------------------------------

	var onceA, onceB sync.Once

	var initB func()
	initA := func() { onceB.Do(initB) }
	initB = func() { onceA.Do(initA) }
	onceA.Do(initA) // deadlock!
}
