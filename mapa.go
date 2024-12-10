package main

import (
	"encoding/json"
	"os"
)

type CamadasJSON struct {
	Dados   []int `json:"data"`
	Largura int   `json:"width"`
	Altura  int   `json:"height"`
}

type MapasJSON struct {
	Ladrilhos []CamadasJSON `json:"layers"`
}

func NovoMapa(caminho string) (*MapasJSON, error) {
	conteúdos, err := os.ReadFile(caminho)
	if err != nil {
		return nil, err
	}

	var mapaJSON MapasJSON

	err = json.Unmarshal(conteúdos, &mapaJSON)
	if err != nil {
		return nil, err
	}

	return &mapaJSON, nil
}
