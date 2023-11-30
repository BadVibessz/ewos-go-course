package main

import (
	"fmt"
	"github.com/ew0s/ewos-to-go-hw/basics1/homework/cell"
)

func main() {
	fmt.Println("\033[31mHello \033[0mWorld") // https://www.shellhacks.com/bash-colors/

	c := cell.CreateCell("станок",
		"станок для дерева",
		"100$",
		"Казань",
		"имеется",
		cell.Row{"😎", "DURA", "DURA"},
		cell.Row{"😎", "AAAAAAAAAAAAAAAAAAAAAAA", "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"},
	)

	//c.Draw(Borderless, Bold, Red)
	c.Draw(cell.Border, cell.Red)
}
