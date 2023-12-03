package cell

import (
	"fmt"
	"slices"
	"strconv"
	"unicode/utf8"

	stringutils "github.com/ew0s/ewos-to-go-hw/basics1/homework/utils/strings"
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

const (
	foregroundConst = 0
	backgroundConst = 0
)

func (c Color) Foreground() Color {
	return c + foregroundConst
}

func (c Color) Background() Color {
	return c + backgroundConst
}

type Mod func(s string) string

func ColorFunc(col Color) Mod {
	return func(s string) string {
		return "\u001B[" + strconv.Itoa(int(col)) + "m" + s + "\u001B[0m"
	}
}

func CharFunc(typ CharType) Mod {
	return func(s string) string {
		return "\u001B[" + strconv.Itoa(int(typ)) + "m" + s + "\u001B[0m"
	}
}

type (
	Row  = [3]string
	Cell []Row
)

func CreateCell(name, desc, price, loc, deliv string, rows ...Row) Cell {
	defaultRowsNum := 5 // anti mnd: Magic number: 5, in <argument> detected (gomnd)

	res := append(make(Cell, 0, len(rows)+defaultRowsNum),
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
	for _, mod := range mods {
		s = mod(s)
	}

	return s
}

func (c *Cell) drawWithBorder(vertBord, horizBord string, upperOff, bottomOff int, mods ...Mod) {
	maxLen := c.maxLenOfRow()

	for i, v := range *c {
		if i == 0 {
			fmt.Println(stringutils.Populate(vertBord, maxLen+upperOff))
		} else {
			fmt.Println(stringutils.Populate(vertBord, maxLen+bottomOff) + horizBord)
		}

		fmt.Print(horizBord)

		line := v[0] + " " + v[1] + ": " + v[2]

		fmt.Print(applyMods(line, mods...))

		fmt.Print(applyMods(stringutils.Populate(" ", maxLen+3-utf8.RuneCountInString(line)), mods...))
		fmt.Println(horizBord)

		fmt.Print(horizBord)

		if i == len(*c)-1 {
			fmt.Println(stringutils.Populate(vertBord, maxLen+bottomOff) + horizBord)
		}
	}
}

func (c *Cell) Draw(typ DrawingType, mods ...Mod) {
	switch typ {
	case Borderless:
		c.drawBorderless(mods...)

	case Border:
		upperOff := 6
		bottomOff := 4
		c.drawWithBorder("_", "|", upperOff, bottomOff, mods...)

	case StarredBorder:
		upperOff := 8
		bottomOff := 4
		c.drawWithBorder("*", "**", upperOff, bottomOff, mods...)

	default:
		c.drawBorderless(mods...)
	}
}
