package cell

import (
	"fmt"
	stringutils "github.com/ew0s/ewos-to-go-hw/basics1/homework/utils/string"
	"slices"
	"strconv"
	"unicode/utf8"
)

type DrawingType int

const (
	Borderless = DrawingType(iota)
	Border
	StarredBorder
)

type CharType int

const (
	Normal     = CharType(0)
	Bold       = CharType(1)
	Underlined = CharType(4)
	Blinking   = CharType(5)
	Reverse    = CharType(7)
)

type Color int

const (
	Black = Color(iota)
	Red
	Green
	Brown
	Blue
	Purple
	Cyan
	LightGray
)

func (c Color) Foreground() Color {
	return c + 30
}

func (c Color) Background() Color {
	return c + 40
}

func ColorFunc(col Color) func(string) string {
	return func(s string) string {
		return "\u001B[" + strconv.Itoa(int(col)) + "m" + s + "\u001B[0m"
	}
}

func CharFunc(typ CharType) func(string) string {
	return func(s string) string {
		return "\u001B[" + strconv.Itoa(int(typ)) + "m" + s + "\u001B[0m"
	}
}

type (
	Row  = [3]string
	Cell []Row

	Mod func(s string) string
)

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

func (c *Cell) maxLenOfRow() int {
	l := make([]int, len(*c))
	for i, v := range *c {
		l[i] = utf8.RuneCountInString(v[0] + " " + v[1] + " " + v[2])
	}

	return slices.Max(l)
}

func (c *Cell) drawBorderless(mods ...Mod) {

	maxLen := c.maxLenOfRow()
	for _, v := range *c {

		line := v[0] + " " + v[1] + ": " + v[2]

		fmt.Print(applyMods(line, mods...))
		fmt.Println(applyMods(stringutils.Populate(" ", maxLen+2-utf8.RuneCountInString(line)), mods...))
	}
}

func applyMods(s string, mods ...Mod) string {
	res := s
	for _, mod := range mods {
		res = mod(res)
	}
	return res
}

func (c *Cell) drawWithBorders(mods ...Mod) {

	maxLen := c.maxLenOfRow()
	for i, v := range *c {

		if i == 0 {
			fmt.Println(stringutils.Populate("_", maxLen+6))
		} else {
			fmt.Println(stringutils.Populate("_", maxLen+4) + "|")
		}

		fmt.Print("|")

		line := v[0] + " " + v[1] + ": " + v[2]

		fmt.Print(applyMods(line, mods...))

		fmt.Print(applyMods(stringutils.Populate(" ", maxLen+3-utf8.RuneCountInString(line)), mods...))
		fmt.Println("|")

		fmt.Print("|")
		if i == len(*c)-1 {
			fmt.Println(stringutils.Populate("_", maxLen+4) + "|")
		}
	}
}

func (c *Cell) drawWithStarredBorder(mods ...Mod) {

	maxLen := c.maxLenOfRow()
	for i, v := range *c {

		if i == 0 {
			fmt.Println(stringutils.Populate("*", maxLen+9))
		} else {
			fmt.Println(stringutils.Populate("*", maxLen+7) + "*")
		}

		fmt.Print("** ")

		line := v[0] + " " + v[1] + ": " + v[2]

		fmt.Print(applyMods(line, mods...))

		fmt.Print(applyMods(stringutils.Populate(" ", maxLen+3-utf8.RuneCountInString(line)), mods...))

		fmt.Println("**")

		fmt.Print("*")
		if i == len(*c)-1 {
			fmt.Println(stringutils.Populate("*", maxLen+7) + "*")
		}
	}

}

func (c *Cell) Draw(typ DrawingType, mods ...Mod) {

	switch typ {
	case Borderless:
		c.drawBorderless(mods...)
	case Border:
		c.drawWithBorders(mods...)
	case StarredBorder:
		c.drawWithStarredBorder(mods...)

	default:
		c.drawBorderless(mods...)
	}
}
