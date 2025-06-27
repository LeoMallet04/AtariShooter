package main

import (
	"log"
	"time"

	"github.com/gdamore/tcell/v2"
)

// extende a estrutura Sprite e o Dir armazena a direcao do disparo

func main() {
	screen, err := tcell.NewScreen()
	if err != nil {
		log.Fatal(err)
	}
	defer screen.Fini()

	if err := screen.Init(); err != nil {
		log.Fatal(err)
	}

	player1 := NewSprite('@', 40, 12)
	player2 := NewSprite('#', 20, 12)
	var bullets []*Bullet
	player1Dir := 'd'
	player2Dir := 'a'

	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	running := true
	for running {
		select {
		case <-ticker.C:
			screen.Clear()
			player1.Draw(screen)
			player2.Draw(screen)

			w, h := screen.Size()
			newBullets := []*Bullet{}
			for _, b := range bullets {
				b.Update()
				if b.X >= 0 && b.Y >= 0 && b.X < w && b.Y < h {
					b.Draw(screen)
					newBullets = append(newBullets, b)
				}
			}
			bullets = newBullets
			screen.Show()

			for screen.HasPendingEvent() {
				ev := screen.PollEvent()
				switch ev := ev.(type) {
				case *tcell.EventKey:
					switch ev.Key() {
					// Player 2 - setas
					case tcell.KeyUp:
						player2Dir = 'w'
						player2.Move('w')
					case tcell.KeyDown:
						player2Dir = 's'
						player2.Move('s')
					case tcell.KeyLeft:
						player2Dir = 'a'
						player2.Move('a')
					case tcell.KeyRight:
						player2Dir = 'd'
						player2.Move('d')
					case tcell.KeyEnter:
						var bx, by int
						switch player2Dir {
						case 'w':
							bx, by = player2.X, player2.Y-1
						case 'a':
							bx, by = player2.X-1, player2.Y
						case 's':
							bx, by = player2.X, player2.Y+1
						case 'd':
							bx, by = player2.X+1, player2.Y
						}
						bullet := NewBullet(bx, by, player2Dir)
						bullet.Char = '*' // Diferencia os tiros do player 2
						bullets = append(bullets, bullet)

					// Player 1 - wasd + e para atirar
					case tcell.KeyRune:
						switch ev.Rune() {
						case 'w', 'a', 's', 'd':
							player1Dir = ev.Rune()
							player1.Move(ev.Rune())
						case 'e':
							var bx, by int
							switch player1Dir {
							case 'w':
								bx, by = player1.X, player1.Y-1
							case 'a':
								bx, by = player1.X-1, player1.Y
							case 's':
								bx, by = player1.X, player1.Y+1
							case 'd':
								bx, by = player1.X+1, player1.Y
							}
							bullets = append(bullets, NewBullet(bx, by, player1Dir))
						case 'q':
							running = false
						}
					case tcell.KeyEsc:
						running = false
					}
				}
			}
		}
	}
}
