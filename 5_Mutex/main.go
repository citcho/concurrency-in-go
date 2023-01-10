package main

import (
	"fmt"
	"sync"
)

func main() {
	// Mutexはクリティカルセクションを保護する手段の一つ
	// クリティカルセクションは共有リソースに対する排他的アクセスを必要とする場所のこと
	// Mutexは「開発者が守らなければいけないメモリに対する同期アクセスの慣習を作る」ことができる
	// 他にも通信を使ったチャネルでのメモリ共有方法もある

	var count int
	var lock sync.Mutex

	increment := func() {
		// lock, unlockで囲んでcountの占有を要求
		// unlockをdeferで記述するイディオムは、panicでも確実に呼び出し、デッドロックを回避するのに役立つ
		lock.Lock()
		defer lock.Unlock()
		count++
		fmt.Printf("Incrementing: %d\n", count)
	}

	decrement := func() {
		lock.Lock()
		defer lock.Unlock()
		count--
		fmt.Printf("Decrementing: %d\n", count)
	}

	var arithmetic sync.WaitGroup

	// インクリメント
	for i := 0; i <= 5; i++ {
		arithmetic.Add(1)
		go func() {
			defer arithmetic.Done()
			increment()
		}()
	}

	// デクリメント
	for i := 0; i <= 5; i++ {
		arithmetic.Add(1)
		go func() {
			defer arithmetic.Done()
			decrement()
		}()
	}

	arithmetic.Wait()
	fmt.Println("Arithmetic complete.")
	// 順不同で出力される
	// Incrementing: 1
	// Decrementing: 0
	// Decrementing: -1
	// Decrementing: -2
	// Decrementing: -3
	// Decrementing: -4
	// Decrementing: -5
	// Incrementing: -4
	// Incrementing: -3
	// Incrementing: -2
	// Incrementing: -1
	// Incrementing: 0
}
