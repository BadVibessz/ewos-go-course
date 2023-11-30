package main

import "fmt"

type (
	Row  = [3]string
	Cell []Row

	DrawingType int

	Mod func(s string) string
)

const (
	Borderless = DrawingType(iota)
	Border     = DrawingType(iota)
)

func Red(s string) string {
	return "\u001B[31m" + s
}

func (c *Cell) drawBorderless(mods ...Mod) {
	for _, v := range *c {

		line := v[0] + " " + v[1] + ": " + v[2]
		for _, mod := range mods {
			line = mod(line)
		}

		fmt.Println(line)
	}
}

func (c *Cell) Draw(typ DrawingType, mods ...Mod) {

	switch typ {
	case Borderless:
		c.drawBorderless(mods...)
	}
}

// todo: обязательные парамы должны быть!
func CreateCell(name, desc, price, loc, deliv string, rows ...Row) Cell {

	res := make(Cell, 0, len(rows)+5)

	res = append(res,
		Row{"💬", "Название", name},
		Row{"📖", "Описание", desc},
		Row{"💵", "Цена", price},
		Row{"📍", "Локация", loc},
		Row{"📦", "Доставка", deliv})

	return append(res, rows...)
}

func (c *Cell) Add(rows ...Row) {
	*c = append(*c, rows...)
}

func main() {
	fmt.Println("\033[31mHello \033[0mWorld") // https://www.shellhacks.com/bash-colors/

	c := CreateCell("станок",
		"станок для дерева",
		"100$",
		"Казань",
		"имеется",
		Row{"😎", "DURA", "DURA"},
	)

	c.Draw(Borderless, Red)

}
