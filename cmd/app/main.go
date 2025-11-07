package main

import (
	"log"
	"sw2p2go/config"
	"sw2p2go/internal/app"
)

// @title           API SW2P2GO - Usuarios
// @version         1.0
// @description     API REST para gestión de usuarios, planes y suscripciones
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      usuario.sw2ficct.lat
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Error cargando configuración: %v", err)
	}

	log.Println("Iniciando servidor...")
	app.Run(cfg)
}
