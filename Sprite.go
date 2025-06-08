package main

import (
	"github.com/gdamore/tcell/v2"
)

type Sprite struct {
	Char rune
	X, Y int
}

func NewSprite(char rune, x, y int) *Sprite {
	return &Sprite{
		Char: char,
		X:    x,
		Y:    y,
	}
}

func (s *Sprite) Draw(screen tcell.Screen) {
	screen.SetContent(
		s.X,
		s.Y,
		s.Char,
		nil,
		tcell.StyleDefault,
	)
}

func (s *Sprite) Move(direction rune) {
	switch direction {
	case 'w':
		s.Y -= 1
	case 'a':
		s.X -= 1
	case 's':
		s.Y += 1
	case 'd':
		s.X += 1
	}

}
