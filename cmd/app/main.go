package main

import (
	"log"
	"sw2p2go/config"
	"sw2p2go/internal/app"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Error cargando configuración: %v", err)
	}

	log.Println("Iniciando servidor...")
	app.Run(cfg)
}
