package main

import (
	"fmt"
	"time"
)

func main() {
	// 関数呼び出しによる起動
	go sayHello()

	// 即時関数による起動
	go func() {
		fmt.Println("hello")
	}()

	// 変数に代入・呼び出し起動
	sayWorld := func() {
		fmt.Println("hello")
	}
	go sayWorld()

	// 合流ポイントがないので実行のタイミングは不確定になる
	// Sleepで処理を待つのはゴルーチンが実行される確率を上げているだけに過ぎない
	time.Sleep(1 * time.Second)
}

func sayHello() {
	fmt.Println("hello")
}
