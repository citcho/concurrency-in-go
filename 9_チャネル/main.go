package main

import (
	"fmt"
	"sync"
)

func main() {
	// チャネルはGoにおける同期処理のプリミティブ
	// 水の流れる川のように、チャネルは情報の流れの水路として機能する
	// メモリに対するアクセスを同期するのに使える一方で、goroutine間の通信に使うのが最適
	// := 演算子を使用して一度にチャネルを生成できるが、個々の手順を分けてみられるのが便利な場合もある
	// 1方向チャネルを初期化することはあまり見られないが、引数や返り値として使用されることが多い

	var c chan string
	c = make(chan string)

	go func(c chan string) <-chan string { // 暗黙的に受信専用チャネルに変換
		c <- "xxx"
		defer close(c) // 受信元でcloseするのではなく、送信元でcloseする
		return c
	}(c)

	// 受信するまで処理がここでブロックされる
	// closedにはチャネルが閉じられたかどうかが真偽値で代入される
	// 閉じられたチャネルからも値は取得できるが、型のデフォルトの値が取得される
	v, closed := <-c
	fmt.Printf("value: %s, closed: %t\n", v, closed) // value: xxx, closed: true

	// ---------------------------------------------------------------------------------------

	// 下流でcloseを行うことでメインGoroutineを終了させるfor range文
	intStream := make(chan any)
	go func() {
		defer close(intStream) // Goroutineを抜ける前にチャネルを閉じる（頻出パターン）
		for i := 0; i <= 5; i++ {
			intStream <- i
		}
	}()

	for integer := range intStream { // チャネルが閉じたら自動でfor文を終了する
		fmt.Printf("%v", integer)
	}
	// 出力
	// 012345

	// ---------------------------------------------------------------------------------------

	// Condのように複数のGoroutineを一度に解放する手段もある
	// 組み合わせやすいので、こういった場合は基本チャネルを採用する
	begin := make(chan any)
	var wg sync.WaitGroup
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			<-begin
			fmt.Printf("%v has begun\n", i)
		}(i)
	}

	fmt.Println("Unblocking goroutines...")
	close(begin)
	wg.Wait()
	// 出力
	// Unblocking goroutines...
	// 3 has begun
	// 2 has begun
	// 0 has begun
	// 1 has begun

	// ---------------------------------------------------------------------------------------

	// 誰がチャネルを所有するのかを明確にするため
	// チャネルの所有者と消費者に分ける
	// 所有者はチャネルの初期化・書き込み・所有権の譲渡・閉じることを扱い、
	// 消費者はチャネルのいつ閉じられたか把握すること・ブロックする操作を扱う
	chanOwner := func() <-chan int {
		resultStream := make(chan int, 5)
		go func() {
			defer close(resultStream)
			for i := 0; i <= 5; i++ {
				resultStream <- i
			}
		}()
		return resultStream
	}

	resultStream := chanOwner()
	for result := range resultStream {
		fmt.Printf("Received: %d\n", result)
	}
	fmt.Println("Done receiving!")
	// 出力
	// Received: 0
	// Received: 1
	// Received: 2
	// Received: 3
	// Received: 4
	// Received: 5
	// Done receiving!
}
