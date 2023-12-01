package main

import (
	"github.com/ew0s/ewos-to-go-hw/basics1/homework/cell"
)

func main() {
	// https://www.shellhacks.com/bash-colors/

	c := cell.CreateCell("станок",
		"станок для дерева",
		"100$",
		"Казань",
		"имеется",
		cell.Row{"😎", "DURA", "DURA"},
		cell.Row{"😎", "AAAAAAAAAAAAAAAAAAAAAAA", "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"},
	)

	c.Draw(cell.Borderless,
		cell.ColorFunc(cell.LightGray.Background()),
		cell.ColorFunc(cell.Purple.Foreground()),
		cell.CharFunc(cell.Bold),
	)

	c.Draw(cell.Border,
		cell.ColorFunc(cell.LightGray.Background()),
		cell.ColorFunc(cell.Purple.Foreground()),
		cell.CharFunc(cell.Bold),
	)

	c.Draw(cell.StarredBorder,
		cell.ColorFunc(cell.LightGray.Background()),
		cell.ColorFunc(cell.Purple.Foreground()),
		cell.CharFunc(cell.Bold),
	)
}
