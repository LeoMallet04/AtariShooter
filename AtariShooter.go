package main

import (
	// "go/printer"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"atari-shooter/PP2PLink"

	"github.com/gdamore/tcell/v2"
)


type GameState struct{
	Bullets []*Bullet
	Players []*Sprite
	Players []*Sprite	
	Players []*Sprite
	bullet    *Bullet	
	bullet    *Bullet	
}


func EncodeGameState(state GameState) string {
	playerStr := SpriteToString(state.Players[0])

	bulletStr := BulletToString(state.bullet)

	return fmt.Sprintf("P[%s]|B[%s]", playerStr,bulletStr)
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
		b, err := BulletFromString(bClean)
		if err != nil {
			return nil, err
		}
			state.Bullets = append(state.Bullets, b)
		b, err := BulletFromString(bClean)
		if err != nil {
			return nil, err
		}
			state.Bullets = append(state.Bullets, b)
		}
	}

	return state, nil
}
		if b.X != -1 || b.Y != -1{ 
			state.Bullets = append(state.Bullets, b)
		}
	}

	return state, nil
}
		if b.X != -1 || b.Y != -1{ 

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
	enderecoLocal := "localhost:8085"
	enderecoRemote := "localhost:8080"

	link := PP2PLink.NewPP2PLink(enderecoLocal, true)


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


	var remoteStateMutex sync.RWMutex
		
		
	var remoteStateMutex sync.RWMutex
		
	localState := &GameState{
		Bullets: []*Bullet{},
		Players: []*Sprite{},
		bullet: NewBullet(-1,-1,'>'),
	}
	
	remoteState := &GameState{
		Bullets: []*Bullet{},
		Players: []*Sprite{},
	}

	localState.Players = initPlayers(1,localState.Players)


	dirs := make([]rune, 2)
	dirs = append(dirs, 'd')

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	running := true
	for running {
		<-ticker.C
		screen.Clear()

		updated := SyncGameState(localState, link, enderecoRemote)
		if updated != nil {
			remoteState = updated
		}
		
		for _, p := range remoteState.Players {
			p.Draw(screen)
		}
		
		for _, p := range localState.Players {
			p.Draw(screen)
		}
		

		w, h := screen.Size() // pega largura (w) e altura (h) da tela
		newBullets := []*Bullet{}



		// atualiza e desenhas as balas
		if localState.bullet.X != -1 || localState.bullet.Y != -1 {
			localState.Bullets = append(localState.Bullets, localState.bullet)
			localState.bullet  = NewBullet(-1, -1, '>')
		}

		for _, b := range localState.Bullets {
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


func initPlayers(playerId int, players []*Sprite) []*Sprite{
	symbol:= rune('@' +playerId)
	x:= 15 *playerId+1
	y:= 10 *playerId+1

	player := NewSprite(playerId,symbol,x,y)

	players = append(players, player)
	
	return players
}


func moveSprites(ev tcell.Event, localState *GameState,dirs[]rune, running *bool){
	if len(localState.Players) == 0 || len(dirs) == 0{
		fmt.Println("Aviso: localState.Players ou dirs está vazio!")
		return
	}
	switch ev := ev.(type) {
	case *tcell.EventKey:
	switch ev.Key() {
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
