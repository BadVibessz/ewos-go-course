package main

import "github.com/ew0s/ewos-to-go-hw/basics1/homework/cell"

// https://www.shellhacks.com/bash-colors/
func main() {
	c := cell.New(
		cell.Row{"💬", "Название", "Станок"},
		cell.Row{"📖", "Описание", "Станок для дерева"},
		cell.Row{"💵", "Цена", "1000 $"},
		cell.Row{"📍", "Локация", "Москва"},
		cell.Row{"📦", "Доставка", "Нет"},
		cell.Row{"🚀", "Опциональная строка 1", "Значение 1"},
		cell.Row{"🎒", "Опциональная строка 2", "Значение 2"},
	)

	c.Draw(c.Borderless(),
		cell.ColorModifier(cell.LightGray.Background()),
		cell.ColorModifier(cell.Purple.Foreground()),
		cell.CharModifier(cell.Bold),
	)

	c.Draw(c.Border(),
		cell.ColorModifier(cell.LightGray.Background()),
		cell.ColorModifier(cell.Purple.Foreground()),
		cell.CharModifier(cell.Bold),
	)

	c.Draw(c.StarredBorder(),
		cell.ColorModifier(cell.LightGray.Background()),
		cell.ColorModifier(cell.Purple.Foreground()),
		cell.CharModifier(cell.Bold),
	)
}
