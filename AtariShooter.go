package main

import (
	"log"
	"time"

	"github.com/gdamore/tcell/v2"
)

// comunica com o outro processo e atualiza as balas
func (newBullets *Bullet[], link PP2PLink) update_procs {

	var receivedBullets []*Bullets
}

func main() {

	//Cria nova tela
	screen, err := tcell.NewScreen()

	//trata o erro
	if err != nil {
		log.Fatal(err)
	}

	//Adia o comportamento de limpar a tela como  ultima coisa
	defer screen.Fini()

	if err := screen.Init(); err != nil {
		log.Fatal(err)
	}

	player := NewSprite('@', 40, 12)
	var bullets []*Bullet
	playerDir := 'd' // direcao inicial do player

	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	running := true
	for running {
		select {
		case <-ticker.C:
			screen.Clear()
			player.Draw(screen)

			w, h := screen.Size() // pega largura (w) e altura (h) da tela

			update_procs(*newBullets)

			// atualiza e desenhas as balas
			newBullets := []*Bullet{}
			for _, b := range bullets {
				b.Update()
				if b.X >= 0 && b.X < w && b.Y >= 0 && b.Y < h {
					b.Draw(screen)
					newBullets = append(newBullets, b)
				}
			}
			bullets = newBullets

			screen.Show()

			newBullets := []*Bullet{}
			for screen.HasPendingEvent() {
				ev := screen.PollEvent()
				switch ev := ev.(type) {
				case *tcell.EventKey:
					switch ev.Rune() {
					case 'w', 'a', 's', 'd':
						playerDir = ev.Rune()
						player.Move(ev.Rune())
					case 'e':
						var bulletX, bulletY int
						switch playerDir {
						case 'w':
							bulletX, bulletY = player.X, player.Y-1
						case 'a':
							bulletX, bulletY = player.X-1, player.Y
						case 's':
							bulletX, bulletY = player.X, player.Y+1
						case 'd':
							bulletX, bulletY = player.X+1, player.Y
						}
						new_Bullets = append(new_Bullets, NewBullet(bulletX, bulletY, playerDir))
					case 'q':
						running = false
					}
				}
			}
		}
	}
}
