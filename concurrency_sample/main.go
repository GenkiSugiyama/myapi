package main

import (
	"fmt"
	"time"
)

func main() {
	boilWater()
	go cutIngredient()

}

func cutIngredient() {
	fmt.Println("start cutting ingredients")

	time.Sleep(1 * time.Second)

	fmt.Println("finish cutting ingredients")
}

func boilWater() {
	fmt.Println("start boiling water")

	time.Sleep(2 * time.Second)

	fmt.Println("finish boiling water")
}
