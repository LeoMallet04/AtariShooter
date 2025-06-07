package main

import (
	"log"

	"github.com/gdamore/tcell/v2"
)

func main(){
	//Cria nova tela
	screen, err := tcell.NewScreen()

	//Trata o erro
	if err != nil{
		log.Fatal(err)
	}

	//Adia o comportamento de limpar a tela como ultima coisa
	defer screen.Fini();

	//Inicializa a tela criada 
	err = screen.Init();
	
	if err != nil {
		log.Fatal(err)
	}

	//game init
	player := NewSprite('@',10,10)

	//Loop principal do jogo
	running := true 
	for running {

		//game init

		screen.Clear()

		player.Draw(screen)
		
		screen.Show()

	

}