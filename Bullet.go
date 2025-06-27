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

func (b *Bullet) ToString() string {
	return fmt.Sprintf("%d;%d;%c", b.Sprite.X, b.Sprite.Y, b.Dir)
}

func BulletFromString(s string) (*Bullet, error) {
	parts := strings.Split(s, ";")
	if len(parts) != 3 {
		return nil, fmt.Errorf("ERRO: conversao de volta para a bala, num args: %s", s)
	}
	x, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, fmt.Errorf("ERRO: conversao de volta para a bala, X: %d", x)
	}
	y, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, fmt.Errorf("ERRO: conversao de volta para a bala, Y: %d", y)
	}
	direcao := []rune(parts[2])
		if len(direcao) != 1 {
		return nil, fmt.Errorf("ERRO: conversao de volta para a bala, direca: %s", string(direcao))
	}
	return &Bullet{Sprite: *NewSprite(direcao[0],x,y), Dir: direcao[0]}, nil
}
