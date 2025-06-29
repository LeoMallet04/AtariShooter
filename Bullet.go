package main

import (
	"fmt"
	"strconv"
	"strings"
)

// extende a estrutura Sprite e o Dir armazena a direcao do disparo
type Bullet struct {
	Sprite
	Dir rune
}

// cria um novo projetil com base na posicao atual do player
func NewBullet(x, y int, dir rune) *Bullet {
	char := map[rune]rune{'w': '^', 'a': '<', 's': 'v', 'd': '>'}[dir]
	return &Bullet{
		Sprite: Sprite{Char: char, X: x, Y: y},
		Dir:    dir,
	}
}

// move a bala
func (b *Bullet) Update() {
	switch b.Dir {
	case 'w':
		b.Y -= 1
	case 'a':
		b.X -= 1
	case 's':
		b.Y += 1
	case 'd':
		b.X += 1
	}
}

func BulletToString(b *Bullet) string {
	return fmt.Sprintf("%c;%d;%d;%c", b.Char, b.X, b.Y, b.Dir)
}

func BulletFromString(s string) (*Bullet, error) {
	parts := strings.Split(s, ";")
	if len(parts) != 4 {
		return nil, fmt.Errorf("ERRO: conversao Bullet, partes: %v", parts)
	}
	char := []rune(parts[0])
	x, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, fmt.Errorf("erro no X: %v", err)
	}
	y, err := strconv.Atoi(parts[2])
	if err != nil {
		return nil, fmt.Errorf("erro no Y: %v", err)
	}
	dir := []rune(parts[3])
	if len(dir) != 1 {
		return nil, fmt.Errorf("erro na direção")
	}
	return &Bullet{Sprite: Sprite{Char: char[0], X: x, Y: y}, Dir: dir[0]}, nil
}