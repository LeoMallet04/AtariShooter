package main

import (
	// "go/printer"
	"log"
	"strings"
	"time"
	// "fmt"

	"atari-shooter/PP2PLink"

	"github.com/gdamore/tcell/v2"
)

// comunica com o outro processo e atualiza as balas
func UpdateBullets(newBullets []*Bullet, playerMove rune, link *PP2PLink.PP2PLink, sendAddress string) {

	// var receivedBullets []*Bullet
	msg := ""
	bulletStr := []string{}
	for _, b := range newBullets {
		bulletStr = append(bulletStr, BulletToString(b))
	}
	msg += strings.Join(bulletStr, ",")

	link.Req <- PP2PLink.PP2PLink_Req_Message{
		To:      sendAddress,
		Message: msg,
	}

	select {
		case recv := <- link.Ind:
			for _, bulletStr := range strings.Split(recv.Message, ",") {
				if bulletStr == ""{ continue }
				bullet, err := BulletFromString(bulletStr)
				if err != nil {
					log.Printf("Erro: %s",err)
				}else{
					newBullets = append(newBullets, bullet)
				}
			}
		case <- time.After(5 * time.Millisecond):
	} 
	
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

			UpadteProcs(*newBullets)

			// atualiza e desenhas as balas
			newBullets := []*Bullet{}
			for _, b := range bullets {
				b.Update()
				//Verifica se a posição x,yu da bala não é nula
				if b.X >= 0 && b.Y >= 0{
					//Verifica se a bala ainda está na tela
					if b.X < w && b.Y < h{
						b.Draw(screen)
						newBullets = append(newBullets, b)
					}
				}
			}
			bullets = newBullets

			screen.Show()

			newBullets := []*Bullet{}
			for screen.HasPendingEvent() {
				ev := screen.PollEvent()
func initPlayers(quantPlayer int, playerC chan Sprite){
	for i:= 0; i < quantPlayer; i++{
		symbol:= rune('@' +i)
		x:= 15 *i+1
		y:= 10 *i+1

		player := NewSprite(symbol,x,y)

		playerC <- *player
	}
	close(playerC)
}


func moveSprites(ev tcell.Event,players[]Sprite,dirs[]rune, bullets *[]*Bullet, running *bool){
	switch ev := ev.(type) {
	case *tcell.EventKey:
	switch ev.Key() {
		// Player 2 - setas
		case tcell.KeyUp:
			dirs[1] = 'w'
			players[1].Move('w')
		case tcell.KeyDown:
			dirs[1] = 's'
			players[1].Move('s')
		case tcell.KeyLeft:
			dirs[1] = 'a'
			players[1].Move('a')
		case tcell.KeyRight:
			dirs[1] = 'd'
			players[1].Move('d')
		case tcell.KeyEnter:
			var bx, by int
			switch dirs[1] {
			case 'w':
				bx, by = players[1].X, players[1].Y-1
			case 'a':
				bx, by = players[1].X-1, players[1].Y
			case 's':
				bx, by = players[1].X, players[1].Y+1
			case 'd':
				bx, by = players[1].X+1, players[1].Y
			}
			bullet := NewBullet(bx, by, dirs[1])
			bullet.Char = '*' // Diferencia os tiros do player 2
			*bullets = append(*bullets, bullet)

		// Player 1 - wasd + e para atirar
		case tcell.KeyRune:
			switch ev.Rune() {
			case 'w', 'a', 's', 'd':
				dirs[0] = ev.Rune()
				players[0].Move(ev.Rune())
			case 'e':
				var bx, by int
				switch dirs[0] {
				case 'w':
					bx, by = players[0].X, players[0].Y-1
				case 'a':
					bx, by = players[0].X-1, players[0].Y
				case 's':
					bx, by = players[0].X, players[0].Y+1
				case 'd':
					bx, by = players[0].X+1, players[0].Y
				}
				*bullets = append(*bullets, NewBullet(bx, by, dirs[0]))
			case 'q':
				*running = false
			}
		case tcell.KeyEsc:
			*running = false
		}
	}
	
}
