package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	salutation := "welcome"
	wg.Add(1)
	go func() {
		defer wg.Done()
		salutation = "hello"
	}()
	wg.Wait()
	fmt.Println(salutation) // hello
	// 変数のコピーに対して操作せず、元の変数を参照している

	for _, lang := range []string{"scala", "javascript", "go"} {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Println(lang)
		}()
	}
	wg.Wait()
	// 出力
	// go go go
	//
	// ゴルーチンは未来の不確定なタイミングにスケジュールされる
	//
	// Goのランタイムは変数langへの参照がまだ保持されているかを知っているので
	// ゴルーチンがメモリにアクセスできるようにヒープに移し、ガベージコレクションしない
	//
	// 要はforが先に終了し、ゴルーチンが実行する時には最後の値である"go"へ参照してしまうということ

	for _, lang := range []string{"scala", "javascript", "go"} {
		wg.Add(1)
		go func(lang string) {
			defer wg.Done()
			fmt.Println(lang)
		}(lang)
	}
	wg.Wait()
	// 出力
	// go scala javascript
	//
	// 正しく動作させるためには値のコピーを引数に渡す必要がある
}
