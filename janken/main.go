package main

import (
	"fmt"
	"go-handson/janken/strategy"
)

func main() {
	player1 := strategy.NewPlayer("Makoto", strategy.NewWinningStrategy())
	player2 := strategy.NewPlayer("Taro", strategy.NewCircularStrategy())

	for i := 0; i < 20; i++ {
		hand1 := player1.NextHand()
		hand2 := player2.NextHand()
		if hand1.IsStrongerThan(hand2) {
			fmt.Printf("Winner: %s\n", player1.ToString())
			player1.Win()
			player2.Lose()
		} else if hand1.IsWeakerThan(hand2) {
			fmt.Printf("Winner: %s\n", player2.ToString())
			player2.Win()
			player1.Lose()
		} else {
			fmt.Println("Even!!!")
		}
	}
	fmt.Println("Finish!")
	fmt.Println(player1.ToString())
	fmt.Println(player2.ToString())
}
