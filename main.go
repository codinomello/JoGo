package main

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct{}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{120, 180, 255, 255})
	ebitenutil.DebugPrint(screen, "Olá, Mundo!")
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ebiten.WindowSize()
}

func main() {
	file, err := os.Open("./img/plamt.png")
	if err != nil {
		log.Fatalf("Erro ao abrir o ícone: %v", err)
	}

	defer file.Close()

	icon, err := png.Decode(file)
	if err != nil {
		log.Fatalf("Erro ao decodificar o ícone: %v", err)
	}

	// ícone da janela
	ebiten.SetWindowIcon([]image.Image{icon})

	// configurações da janela
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("JoGo")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	// execução do jogo
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatalf("Erro ao executar o jogo: %v", err)
	}
}
