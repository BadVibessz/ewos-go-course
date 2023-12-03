package main

import "github.com/ew0s/ewos-to-go-hw/basics1/homework/cell"

// https://www.shellhacks.com/bash-colors/
func main() {
	c := cell.CreateCell("Название",
		"Описание",
		"1000 $",
		"Москва",
		"Нет",
		cell.Row{"🚀", "Опциональная строка 1", "Значение 1"},
		cell.Row{"🎒", "Опциональная строка 2", "Значение 2"},
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
