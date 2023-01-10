package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	// sync.WaitGroupは、ひとまとまりの並行処理があった時
	// ・その結果を気にしない
	// ・他に結果を収集する手段がある
	// といった場合、並行処理の完了を待つ手段として有効
	// どちらの前提も当てはまらない場合はselect文を使う方が良い
	var wg1 sync.WaitGroup

	wg1.Add(1)
	go func() {
		defer wg1.Done()
		fmt.Println("1st goroutine sleeping...")
		time.Sleep(1 * time.Second)
	}()

	wg1.Add(1)
	go func() {
		defer wg1.Done()
		fmt.Println("2nd goroutine sleeping...")
		time.Sleep(2 * time.Second)
	}()

	wg1.Wait()
	fmt.Println("All goroutine complete.")
	// 2nd goroutine sleeping...
	// 1st goroutine sleeping...
	// All goroutine complete.

	// Addの呼び出しは監視対象のgoroutineの前に書くのが慣習
	// 以下は一度に監視する書き方
	hello := func(wg *sync.WaitGroup, id int) {
		defer wg.Done()
		fmt.Printf("Hello from %v!\n", id)
	}

	const numGreeters = 5
	var wg2 sync.WaitGroup

	wg2.Add(numGreeters)
	for i := 0; i < numGreeters; i++ {
		go hello(&wg2, i+1)
	}
	wg2.Wait()
	// 順不同で出力される
	// Hello from 5!
	// Hello from 4!
	// Hello from 3!
	// Hello from 2!
	// Hello from 1!
}
