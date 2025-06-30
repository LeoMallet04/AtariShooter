package main

import (
	// "go/printer"
	"fmt"
	"log"
	"strings"
	"time"

	"atari-shooter/PP2PLink"

	"github.com/gdamore/tcell/v2"
)


type GameState struct{
	Bullets []*Bullet
	Players []*Sprite	
}


func EncodeGameState(state GameState) string {
	playerStr := SpriteToString(state.Players[0])

	bulletStrs := []string{}
	
	for _, b := range state.Bullets{
		bulletStrs = append(bulletStrs, BulletToString(b))
	}
	var bulletStrings = strings.Join(bulletStrs, ",")


	return fmt.Sprintf("P[%s]|B[%s]", playerStr,bulletStrings)
}

func DecodeGameState(s string) (*GameState, error){
	parts := strings.Split(s,"|")

	if len(parts) != 2 {
		log.Printf("Mensagem de Erro: %s",s)
		return nil, fmt.Errorf("FORMATO INVÁLIDO")
	}

	pClean := strings.TrimPrefix(strings.TrimSuffix(parts[0], "]"), "P[")
	bClean := strings.TrimPrefix(strings.TrimSuffix(parts[1], "]"), "B[")

    state := &GameState{}


	if pClean != "" {
		p, err := SpriteFromString(pClean)
		if err != nil {
			return nil, err
		}
		state.Players = append(state.Players, p)
	}
	
	if bClean != ""{
		for _, bs := range strings.Split(bClean, ","){
			b, err := BulletFromString(bs)
			if err != nil {
				return nil, err
			}
			state.Bullets = append(state.Bullets, b)
		}
	}

	return state, nil
}

// comunica com o outro processo e atualiza as balas

func SyncGameState(localState *GameState, link *PP2PLink.PP2PLink, sendAddress string) *GameState{
	msg := EncodeGameState(*localState)

	link.Req <- PP2PLink.PP2PLink_Req_Message{
		To: sendAddress,
		Message: msg,
	}

	select {
	case recv := <-link.Ind:
		if recv.Message == "" {
			return nil
		}
		remoteState, err := DecodeGameState(recv.Message)
		if err != nil {
			log.Printf("Erro ao decodificar: '%s' - erro: %v", recv.Message, err)
			return nil
		}
		// log.Println("Mensagem recebida:", recv.Message)
		return remoteState
	case <-time.After(100 * time.Millisecond):
		return nil
	}
}


func main() {
	enderecoLocal := "192.168.1.68:8080"
	enderecoRemote := "192.168.1.68:8085"

	link := PP2PLink.NewPP2PLink(enderecoLocal, true)

	link.Start(enderecoRemote)

	playerC:= make(chan Sprite)


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

	localState := &GameState{
		Bullets: []*Bullet{},
		Players: []*Sprite{},
	}

	remoteState := &GameState{
		Bullets: []*Bullet{},
		Players: []*Sprite{},
	}

	

	players:= []Sprite{}
	dirs:= []rune{}
	for p := range playerC{
		players =append(players, p)
		dirs = append(dirs, 'd')	
		localState.Players = append(localState.Players, &players[len(players)-1])
	}

	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	running := true
	for running {
		<-ticker.C
		screen.Clear()

		for _,p:= range append(localState.Players, remoteState.Players...){
			p.Draw(screen)
		}

		w, h := screen.Size() // pega largura (w) e altura (h) da tela
		recvState := SyncGameState(localState,link,enderecoRemote)
		if recvState != nil{
			remoteState = recvState
		}

		// atualiza e desenhas as balas
		allBullets := append(localState.Bullets,remoteState.Bullets...)
		newBullets := []*Bullet{}

		for _, b := range allBullets {
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
		localState.Bullets = newBullets

		screen.Show()
		for screen.HasPendingEvent() {
			ev := screen.PollEvent()
			moveSprites(ev, localState, dirs, &running)
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


func moveSprites(ev tcell.Event, localState *GameState,dirs[]rune, running *bool){
	switch ev := ev.(type) {
	case *tcell.EventKey:
	switch ev.Key() {
		// Player 2 - setas
		case tcell.KeyUp:
			dirs[1] = 'w'
			localState.Players[1].Move('w')
		case tcell.KeyDown:
			dirs[1] = 's'
			localState.Players[1].Move('s')
		case tcell.KeyLeft:
			dirs[1] = 'a'
			localState.Players[1].Move('a')
		case tcell.KeyRight:
			dirs[1] = 'd'
			localState.Players[1].Move('d')
		case tcell.KeyEnter:
			var bx, by int
			switch dirs[1] {
			case 'w':
				bx, by = localState.Players[1].X, localState.Players[1].Y-1
			case 'a':
				bx, by = localState.Players[1].X-1, localState.Players[1].Y
			case 's':
				bx, by = localState.Players[1].X, localState.Players[1].Y+1
			case 'd':
				bx, by = localState.Players[1].X+1, localState.Players[1].Y
			}
			bullet := NewBullet(bx, by, dirs[1])
			bullet.Char = '*' // Diferencia os tiros do player 2
			localState.Bullets = append(localState.Bullets, bullet)

		// Player 1 - wasd + e para atirar
		case tcell.KeyRune:
			switch ev.Rune() {
			case 'w', 'a', 's', 'd':
				dirs[0] = ev.Rune()
				localState.Players[0].Move(ev.Rune())
			case 'e':
				var bx, by int
				switch dirs[0] {
				case 'w':
					bx, by = localState.Players[0].X, localState.Players[0].Y-1
				case 'a':
					bx, by = localState.Players[0].X-1, localState.Players[0].Y
				case 's':
					bx, by = localState.Players[0].X, localState.Players[0].Y+1
				case 'd':
					bx, by = localState.Players[0].X+1, localState.Players[0].Y
				}
				bullet := NewBullet(bx,by,dirs[0])
				localState.Bullets = append(localState.Bullets, bullet)
			case 'q':
				*running = false
			}
		case tcell.KeyEsc:
			*running = false
		}
	}
	
}
