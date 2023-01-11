package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	// Condはなんらかのシグナルを受けるまで、メインgoroutineを待機させたりする時に利用できる
	// チャネルの受信を待つ方法もあるが、Condを使うべき時は
	// ①大前提： 非同期・並行処理が発生する。
	// ②ある変数の状態の遷移（イベント）が起こった事実を「繰り返し」確認する必要がある。
	// ③その状態変化が起こるまでは確認スレッドは待機している必要がある。
	// ④ただし、確認するタイミングは確認スレッドから見て外部要因によって決まる（例：カウントがある一定の数値を超えた、バッファが埋まった、等）。
	// ⑤上記の処理を繰り返し行う必要がある
	// Mutexと同じで、Lock,Unlockは読み込み・書き込みの時に記述する

	c := sync.NewCond(&sync.Mutex{})
	queue := make([]any, 0, 10)

	removeFromQueue := func(delay time.Duration) {
		time.Sleep(delay)
		c.L.Lock()
		queue = queue[1:]
		fmt.Println("Removed from queue")
		c.L.Unlock()
		c.Signal()
	}

	for i := 0; i < 10; i++ {
		c.L.Lock()
		for len(queue) == 2 {
			c.Wait()
			fmt.Println("Done wait")
		}
		fmt.Println("Adding to queue")
		queue = append(queue, struct{}{})
		go removeFromQueue(1 * time.Second)
		c.L.Unlock()
	}
	// 出力
	// Adding to queue
	// Adding to queue
	// Removed from queue
	// Done wait
	// Adding to queue
	// Removed from queue
	// Done wait
	// Adding to queue
	// Removed from queue
	// Done wait
	// Adding to queue
	// Removed from queue
	// Done wait
	// Adding to queue
	// Removed from queue
	// Done wait
	// Adding to queue
	// Removed from queue
	// Done wait
	// Adding to queue
	// Removed from queue
	// Removed from queue
	// Done wait
	// Adding to queue
	// Adding to queue

	// Broadcast()を使用した例

	// Clickedという状態（Cond）を作成
	type Button struct {
		Clicked *sync.Cond
	}
	button := Button{Clicked: sync.NewCond(&sync.Mutex{})}

	subscribe := func(c *sync.Cond, fn func()) {
		var goroutineRunning sync.WaitGroup
		goroutineRunning.Add(1)
		go func() {
			goroutineRunning.Done()
			c.L.Lock()
			defer c.L.Unlock()
			c.Wait()
			fn()
		}()
		goroutineRunning.Wait()
	}

	var clickRegistered sync.WaitGroup
	clickRegistered.Add(3)
	subscribe(button.Clicked, func() {
		fmt.Println("Maximizing window.")
		clickRegistered.Done()
	})
	subscribe(button.Clicked, func() {
		fmt.Println("Displaying annoying dialog box!")
		clickRegistered.Done()
	})
	subscribe(button.Clicked, func() {
		fmt.Println("Mouse clicked.")
		clickRegistered.Done()
	})

	button.Clicked.Broadcast()

	clickRegistered.Wait()
	// 出力
	// Mouse clicked.
	// Maximizing window.
	// Displaying annoying dialog box!
}
