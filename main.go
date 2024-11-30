package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Sprite struct {
	Imagem *ebiten.Image
	X, Y   float64
}

type Jogador struct {
	*Sprite
	Vida uint
}

type Inimigo struct {
	*Sprite
	Seguir bool
}

type Animal struct {
	*Sprite
	Seguir bool
}

type Item struct {
	*Sprite
	Poção uint
}

type Jogo struct {
	Player   *Jogador
	Inimigos []*Inimigo
	Animais  []*Animal
	Itens    []*Item
}

func (g *Jogo) Update() error {
	// input do teclado
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.Player.X += 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.Player.X -= 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		g.Player.Y -= 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		g.Player.Y += 2
	}

	for _, sprite := range g.Inimigos {
		if sprite.Seguir {
			if sprite.X < g.Player.X {
				sprite.X += 1
			} else if sprite.X > g.Player.X {
				sprite.X -= 1
			}
			if sprite.Y < g.Player.Y {
				sprite.Y += 1
			} else if sprite.Y > g.Player.Y {
				sprite.Y -= 1
			}
		}
	}

	for _, item := range g.Itens {
		if g.Player.X == item.X {
			g.Player.Vida += item.Poção
			fmt.Printf("O baú foi aberto! Vida: %d\n", g.Player.Vida)
		}
	}

	return nil
}

func (g *Jogo) Draw(screen *ebiten.Image) {
	largura, altura := ebiten.WindowSize()
	janela := ebiten.NewImage(largura, altura)

	// preenchimento da janela
	janela.Fill(color.RGBA{120, 180, 255, 255})

	// desenho do jogador
	opções := ebiten.DrawImageOptions{}
	opções.GeoM.Translate(g.Player.X, g.Player.Y)

	janela.DrawImage(
		g.Player.Imagem.SubImage(
			image.Rect(0, 0, 32, 32),
		).(*ebiten.Image),
		&opções,
	)
	opções.GeoM.Reset()

	// desenho de inimigos
	for _, sprite := range g.Inimigos {
		opções.GeoM.Translate(sprite.X, sprite.Y)

		janela.DrawImage(
			sprite.Imagem.SubImage(
				image.Rect(0, 0, 32, 32),
			).(*ebiten.Image),
			&opções,
		)
		opções.GeoM.Reset()
	}

	// desenho de animais
	for _, sprite := range g.Animais {
		opções.GeoM.Translate(sprite.X, sprite.Y)

		janela.DrawImage(
			sprite.Imagem.SubImage(
				image.Rect(0, 0, 32, 32),
			).(*ebiten.Image),
			&opções,
		)
		opções.GeoM.Reset()
	}

	// desenho de itens
	for _, sprite := range g.Itens {
		opções.GeoM.Translate(sprite.X, sprite.Y)

		janela.DrawImage(
			sprite.Imagem.SubImage(
				image.Rect(0, 0, 32, 32),
			).(*ebiten.Image),
			&opções,
		)
		opções.GeoM.Reset()
	}

	escala := &ebiten.DrawImageOptions{}
	escala.GeoM.Scale(1.0, 1.0)
	screen.DrawImage(janela, escala)
}

func (g *Jogo) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	// atribuição do ícone
	file, err := os.Open("img/plamt.png")
	if err != nil {
		log.Fatalf("Erro ao abrir o ícone: %v", err)
	}

	defer file.Close()

	ícone, err := png.Decode(file)
	if err != nil {
		log.Fatalf("Erro ao decodificar o ícone: %v", err)
	}

	// ícone da janela
	ebiten.SetWindowIcon([]image.Image{ícone})

	// configurações da janela
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("JoGo")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	// configuração do jogador
	player, _, err := ebitenutil.NewImageFromFile("assets/jogador/jogador.png")
	if err != nil {
		log.Fatalf("Erro ao configurar o jogador: %v", err)
	}

	// configuração do esqueleto
	esqueleto, _, err := ebitenutil.NewImageFromFile("assets/inimigos/esqueleto.png")
	if err != nil {
		log.Fatalf("Erro ao configurar o esqueleto: %v", err)
	}

	// configuração do slime
	slime, _, err := ebitenutil.NewImageFromFile("assets/inimigos/slime-roxo.png")
	if err != nil {
		log.Fatalf("Erro ao configurar o slime roxo: %v", err)
	}

	// configuração da galinha
	galinha, _, err := ebitenutil.NewImageFromFile("assets/animais/galinha/galinha.png")
	if err != nil {
		log.Fatalf("Erro ao configurar a galinha: %v", err)
	}

	// configuração do baú
	baú, _, err := ebitenutil.NewImageFromFile("assets/decoração/baú.png")
	if err != nil {
		log.Fatalf("Erro ao configurar o baú: %v", err)
	}

	jogo := Jogo{
		// jogador
		Player: &Jogador{
			Sprite: &Sprite{
				Imagem: player,
				X:      100.0,
				Y:      100.0,
			},
			Vida: 100,
		},

		// inimigos
		Inimigos: []*Inimigo{
			{
				&Sprite{
					Imagem: esqueleto,
					X:      150.0,
					Y:      120.0,
				},
				true,
			},
			{
				&Sprite{
					Imagem: slime,
					X:      130.0,
					Y:      170.0,
				},
				false,
			},
		},

		// animais
		Animais: []*Animal{
			{
				&Sprite{
					Imagem: galinha,
					X:      50.0,
					Y:      120.0,
				},
				false,
			},
		},

		// itens
		Itens: []*Item{
			{
				&Sprite{
					Imagem: baú,
					X:      200.0,
					Y:      100.0,
				},
				1,
			},
		},
	}

	// execução do jogo
	if err := ebiten.RunGame(&jogo); err != nil {
		log.Panicf("Erro ao executar o jogo: %v", err)
	}
}
