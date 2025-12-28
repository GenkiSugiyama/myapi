package main

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

var wg sync.WaitGroup

func main() {
	// メインゴルーチンは別ゴルーチンで動く処理を待たずに終わることがある
	// 別ゴルーチンで動く処理を待つにはsync.WaitGroupを用意する
	// wg.Addで待つゴルーチン処理の数を指定する、追加された数分待ちがカウントアップされる
	// wg.Waitは待ちが0になったらメインを終了する？
	// 別ゴルーチン側ではwg.Doneで待ちをカウントダウンする
	// wg.Add(2)
	// go cutIngredient()
	// go boilWater()
	// wg.Wait()

	// チャネルを使った並行処理の実装
	// チャネルの初期値はnilなので、var ch chan int のような宣言だとnilになってしまう
	// チャネルを使用するには組み込み関数のmakeで作る
	ch1, ch2 := make(chan int), make(chan string)
	// チャネルは使い終わったらclose()が必要
	defer close(ch1)
	defer close(ch2)

	go doubleInt(4, ch1)
	go doubleString("hello", ch2)

	// チャネルはチャネルへの送信・チャネルからの受信時にチャネル側の準備ができていなかったら自動的に待ち状態になる
	for i := 0; i < 2; i++ {
		select {
		case numResult := <-ch1:
			fmt.Println(numResult)
		case strResult := <-ch2:
			fmt.Println(strResult)
		}
	}
}

func doubleInt(src int, intCh chan<- int) {
	// 第二引数のチャネルは送信専用なのでdoubleInt関数ではintChには値を送信しかできない
	// wg.WaitGroupを使わなくても待ち合わせが行われる
	result := src * 2
	intCh <- result
}

func doubleString(src string, strChan chan<- string) {
	result := strings.Repeat(src, 2)
	strChan <- result
}

func cutIngredient() {
	defer wg.Done()
	fmt.Println("start cutting ingredients")

	time.Sleep(1 * time.Second)

	fmt.Println("finish cutting ingredients")
}

func boilWater() {
	defer wg.Done()
	fmt.Println("start boiling water")

	time.Sleep(2 * time.Second)

	fmt.Println("finish boiling water")
}
