package cell

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
	"unicode/utf8"
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
	foregroundConst = 30
	backgroundConst = 40
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
		return fmt.Sprintf("\u001B[%vm%v\u001B[0m", strconv.Itoa(int(col)), s)
	}
}

func CharFunc(typ CharType) Mod {
	return func(s string) string {
		return fmt.Sprintf("\u001B[%vm%v\u001B[0m", strconv.Itoa(int(typ)), s)
	}
}

const numOfValues = 3

type (
	Row  = [numOfValues]string
	Cell []Row
)

// GetRequiredRowNames exported because user need to know what's required (or it's better to store array like this in global var?)
func GetRequiredRowNames() []string {
	return []string{"Название", "Описание", "Цена", "Локация", "Доставка"}
}

func (c *Cell) validate() bool {
	for _, name := range GetRequiredRowNames() {
		ind := slices.IndexFunc(*c, func(r Row) bool { return r[1] == name })
		if ind == -1 {
			return false
		}
	}

	return true
}

func New(rows ...Row) *Cell {
	cell := Cell(rows)

	valid := cell.validate()
	if !valid {
		fmt.Println("Invalid cell")
		return nil
	}

	return &cell
}

func (c *Cell) maxLenOfRow() int {
	l := make([]int, len(*c))

	for i, row := range *c {
		l[i] += utf8.RuneCountInString(row[0]) + utf8.RuneCountInString(row[1]) + utf8.RuneCountInString(row[2])
	}

	numOfSpaces := 2

	return slices.Max(l) + numOfSpaces
}

func (c *Cell) drawBorderless(mods ...Mod) {
	maxLen := c.maxLenOfRow()
	off := 2

	for _, row := range *c {
		line := fmt.Sprintf("%s %s: %s", row[0], row[1], row[2])

		fmt.Print(applyMods(line, mods...))

		spacesCount := maxLen + off - utf8.RuneCountInString(line)
		fmt.Println(applyMods(strings.Repeat(" ", spacesCount), mods...))
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

	for i, row := range *c {
		if i == 0 {
			fmt.Println(strings.Repeat(vertBord, maxLen+upperOff))
		} else {
			fmt.Println(strings.Repeat(vertBord, maxLen+bottomOff) + horizBord)
		}

		fmt.Print(horizBord)

		line := fmt.Sprintf("%s %s: %s", row[0], row[1], row[2])

		fmt.Print(applyMods(line, mods...))

		fmt.Print(applyMods(strings.Repeat(" ", maxLen+3-utf8.RuneCountInString(line)), mods...))
		fmt.Println(horizBord)

		fmt.Print(horizBord)

		if i == len(*c)-1 {
			fmt.Println(strings.Repeat(vertBord, maxLen+bottomOff) + horizBord)
		}
	}
}

type DrawingFunc func(mods ...Mod)

func (c *Cell) Borderless() DrawingFunc {
	return c.drawBorderless
}

func (c *Cell) Border() DrawingFunc {
	return func(mods ...Mod) {
		upperOff := 6
		bottomOff := 4
		c.drawWithBorder("_", "|", upperOff, bottomOff, mods...)
	}
}

func (c *Cell) StarredBorder() DrawingFunc {
	return func(mods ...Mod) {
		upperOff := 8
		bottomOff := 4
		c.drawWithBorder("*", "**", upperOff, bottomOff, mods...)
	}
}

func (c *Cell) Draw(drawingFunc DrawingFunc, mods ...Mod) {
	if c == nil {
		fmt.Println("Cell is nil")
		return
	}

	drawingFunc(mods...)
}
