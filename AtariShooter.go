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

type GameState struct{
	Bullets []*Bullet
	Players []*Sprite
}


func EncodeGameState(state GameState) string {
	playerStrs := []string{}

	for _, p := range state.Players{
		playerStrs = append(playerStrs, SpriteToString(p)) 
	}

	bulletStrs := []string{}
	
	for _, b := range state.Bullets{
		bulletStrs = append(bulletStrs, BulletToString(b))
	}

	var playerStrings = strings.Join(playerStrs, ",")
	var bulletStrings = strings.Join(bulletStrs, ",")


	return fmt.Sprintf("P[%s]|B[%s]", playerStrings,bulletStrings)
}
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
	enderecoLocal := "localhost:8080"
	enderecoRemote := "localhost:8085"

	link := PP2PLink.NewPP2PLink(enderecoLocal, true)

	// link.Req <- PP2PLink.PP2PLink_Req_Message{
	// 	To: enderecoRemote,
	// 	Message: "Oi, tudo bem?",
	// }

	// for i := 0; i < 5; i++{
	// 	msg := <- link.Ind
	// 	fmt.Printf("%s Tudo simm\n", msg)
	// }

	playerC:= make(chan Sprite)
	// bulletC:= make(chan Bullet)

	go initPlayers(2,playerC)

	//Cria nova tela
	screen, err := tcell.NewScreen()
	if err != nil {
		log.Fatal(err)
	}

	//Adia o comportamento de limpar a tela como  ultima coisa
	defer screen.Fini()

	if err := screen.Init(); err != nil {
		log.Fatal(err)
	}

	players:= []Sprite{}
	dirs:= []rune{}
	for p := range playerC{
		players =append(players, p)
		dirs = append(dirs, 'd')
	}

	var bullets = []*Bullet{}
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	running := true
	for running {
		<-ticker.C
		screen.Clear()
		for _,p:= range players{
			p.Draw(screen)
		}

		w, h := screen.Size() // pega largura (w) e altura (h) da tela

		UpdateBullets(bullets,players[0].Char,link,enderecoRemote)

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

		// newBullets := []*Bullet{}
		for screen.HasPendingEvent() {
			ev := screen.PollEvent()
			moveSprites(ev,players,dirs,&bullets,&running)
		}
	}
}



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
