package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
)

type Sprite struct {
	Id int
	Char rune
	X, Y int
}

func NewSprite(id int, char rune, x, y int) *Sprite {
	return &Sprite{
		Id: id,
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

func SpriteToString(s *Sprite) string{
	return fmt.Sprintf("%d;%c;%d;%d",s.Id,s.Char, s.X, s.Y)
}

func SpriteFromString(str string) (*Sprite, error){
	parts := strings.Split(str, ";")
	if len(parts) != 4 {
		return nil, fmt.Errorf("ERRO: conversao de volta para a bala, número de args: %s", str)
	}
	id, err := strconv.Atoi(parts[0])
	if err != nil{
		return nil, fmt.Errorf("ERRO: conversão de volta para a bala, Id: %d",id)
	}
	char := []rune(parts[1])
	if len(char) != 1{
		return nil, fmt.Errorf("ERRO: conversão de volta para a bala, Char: %c",char)
	}
	x, err := strconv.Atoi(parts[2])
	if err != nil {
		return nil, fmt.Errorf("ERRO: conversao de volta para a bala, X: %d", x)
	}
	y, err := strconv.Atoi(parts[3])
	if err != nil {
		return nil, fmt.Errorf("ERRO: conversao de volta para a bala, Y: %d", y)
	}
	return &Sprite{Id: id, Char: char[0], X: x, Y: y}, nil
}
