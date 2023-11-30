package cell

import (
	"fmt"
	stringutils "github.com/ew0s/ewos-to-go-hw/basics1/homework/utils/string"
	"slices"
	"unicode/utf8"
)

//*******************
//**  name * desk **
//*******************

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

// todo: more modes
func Red(s string) string {
	return "\u001B[31m" + s + "\u001B[0m"
}

func Bold(s string) string {
	return "\u001B[1m " + s + "\u001B[0m"
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

func (c *Cell) drawWithBorders(mods ...Mod) {

	// todo: encapsulate
	l := make([]int, len(*c))
	for i, v := range *c {
		l[i] = utf8.RuneCountInString(v[0] + " " + v[1] + " " + v[2])
	}

	maxLen := slices.Max(l)

	for i, v := range *c {

		if i == 0 {
			fmt.Println(stringutils.Populate("_", maxLen+6))
		} else {
			fmt.Println(stringutils.Populate("_", maxLen+4) + "|")
		}

		fmt.Print("|")

		line := v[0] + " " + v[1] + ": " + v[2]

		spaces := ""
		for k := 0; k < maxLen+3-utf8.RuneCountInString(line); k++ {
			spaces += " "
		}

		for _, mod := range mods {
			line = mod(line)
		}

		fmt.Print(line)

		fmt.Print(spaces)
		fmt.Println("|")

		fmt.Print("|")
		//verticalBord = ""
		//for k := 0; k < maxLen-2; k++ {
		//	verticalBord += "_"
		//}

		if i == len(*c)-1 {
			fmt.Println(stringutils.Populate("_", maxLen+4) + "|")
		}
	}
}

func (c *Cell) Draw(typ DrawingType, mods ...Mod) {

	switch typ {
	case Borderless:
		c.drawBorderless(mods...)
	case Border:
		c.drawWithBorders(mods...)
	}
}

func CreateCell(name, desc, price, loc, deliv string, rows ...Row) Cell {

	res := append(make(Cell, 0, len(rows)+5),
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
