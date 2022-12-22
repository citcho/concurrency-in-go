package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	sayHello := func() {
		defer wg.Done()
		fmt.Println("hello")
	}

	wg.Add(1)

	go sayHello()

	// 合流ポイントを作成し、メインゴルーチンをブロックする
	wg.Wait()
}
