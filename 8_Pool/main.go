package main

import (
	"fmt"
	"sync"
)

func main() {
	// sync.Poolはオブジェクトプールパターンを並行処理で安全な形で実装したもの
	// 使うものを決まった数だけ作る方法
	// コストの高いデータベースコネクションを再利用する時などに使用される
	pool := &sync.Pool{
		New: func() any {
			fmt.Println("creating new instance.")
			return struct{}{}
		},
	}

	pool.Get()             // poolの中に空の構造体が無いのでNewを実行
	instance := pool.Get() // poolの中に空の構造体が無いのでNewを実行
	pool.Put(instance)     // 空の構造体をpoolに戻す
	pool.Get()             // poolに戻された空の構造体を再利用
	// 出力
	// creating new instance.
	// creating new instance.
}
