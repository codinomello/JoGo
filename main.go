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
	Imagem  *ebiten.Image
	X, Y    float64
	Display bool
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
	MapaJSON *MapasJSON
	MapaImg  *ebiten.Image
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
	janela.Fill(color.RGBA{50, 135, 45, 1})

	// desenho do jogador
	opções := ebiten.DrawImageOptions{}

	// loop das camadas do mapa
	for _, camadas := range g.MapaJSON.Ladrilhos {
		// loop dos ladrilhos dos dados de cada camada
		for index, id := range camadas.Dados {
			// posição de cada ladrilho
			x := index % camadas.Largura
			y := index / camadas.Altura

			// converte a posição para pixels
			x *= 16
			y *= 16

			// pega a posição na imagem onde o id se localiza
			iX := (id - 1) % 3
			iY := (id - 1) / 3

			// converte a posição para pixels
			iX *= 16
			iY *= 16

			// seta a função para desenhar os ladrilhos em x e y
			opções.GeoM.Translate(float64(x), float64(y))

			// desenha os ladrilhos
			janela.DrawImage(
				// pega o exato ladrilho dos assets
				g.MapaImg.SubImage(image.Rect(iX, iY, iX+16, iY+16)).(*ebiten.Image),
				&opções,
			)
			// reseta o loop e vai para o próximo ladrilho
			opções.GeoM.Reset()
		}
	}

	// seta a transformação dos vértices do jogador
	opções.GeoM.Translate(g.Player.X, g.Player.Y)

	// desenho do player
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
	arquivo, err := os.Open("icons/gopher.png")
	if err != nil {
		log.Fatalf("Erro ao abrir o ícone: %v", err)
	}

	defer arquivo.Close()

	ícone, err := png.Decode(arquivo)
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
	slime, _, err := ebitenutil.NewImageFromFile("assets/inimigos/slime.png")
	if err != nil {
		log.Fatalf("Erro ao configurar o slime: %v", err)
	}

	// configuração da gosma
	gosma, _, err := ebitenutil.NewImageFromFile("assets/inimigos/gosma.png")
	if err != nil {
		log.Fatalf("Erro ao configurar a gosma: %v", err)
	}

	// configuração da galinha
	galinha, _, err := ebitenutil.NewImageFromFile("assets/animais/galinha/galinha.png")
	if err != nil {
		log.Fatalf("Erro ao configurar a galinha: %v", err)
	}
	// configuração da vaca
	vaca, _, err := ebitenutil.NewImageFromFile("assets/animais/vaca/vaca.png")
	if err != nil {
		log.Fatalf("Erro ao configurar a vaca: %v", err)
	}
	// configuração do porco
	porco, _, err := ebitenutil.NewImageFromFile("assets/animais/porco/porco.png")
	if err != nil {
		log.Fatalf("Erro ao configurar o porco: %v", err)
	}

	// configuração do baú
	baú, _, err := ebitenutil.NewImageFromFile("assets/decoração/baú.png")
	if err != nil {
		log.Fatalf("Erro ao configurar o baú: %v", err)
	}

	// configuração dos ladrilhos do mapa
	mapaImg, _, err := ebitenutil.NewImageFromFile("assets/ladrilhos/gramas.png")
	if err != nil {
		log.Fatalf("Erro ao carregar os ladrilhos: %v", err)
	}

	// carregamento do mapa
	mapaJSON, err := NovoMapa("assets/mapas/mapa.json")
	if err != nil {
		log.Panicf("Erro ao carregar o mapa: %v", err)
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
					X:      200.0,
					Y:      120.0,
				},
				false,
			},
			{
				&Sprite{
					Imagem: slime,
					X:      190.0,
					Y:      30.0,
				},
				false,
			},
			{
				&Sprite{
					Imagem: gosma,
					X:      220.0,
					Y:      180.0,
				},
				false,
			},
		},

		// animais
		Animais: []*Animal{
			{
				&Sprite{
					Imagem: galinha,
					X:      130.0,
					Y:      170.0,
				},
				false,
			},
			{
				&Sprite{
					Imagem: vaca,
					X:      50.0,
					Y:      120.0,
				},
				false,
			},
			{
				&Sprite{
					Imagem: porco,
					X:      10.0,
					Y:      190.0,
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
		MapaJSON: mapaJSON,
		MapaImg:  mapaImg,
	}

	// execução do jogo
	if err := ebiten.RunGame(&jogo); err != nil {
		log.Panicf("Erro ao executar o jogo: %v", err)
	}
}
