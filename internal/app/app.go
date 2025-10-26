package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sw2p2go/config"
	v1 "sw2p2go/internal/controller/http/v1"
	"sw2p2go/internal/middleware"
	"sw2p2go/internal/usecase/repositories"
	"sw2p2go/internal/usecase/services"
	"syscall"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type App struct {
	config   *config.Config
	router   *v1.Router
	database *mongo.Database
}

func NewApp(cfg *config.Config) *App {
	return &App{
		config: cfg,
	}
}

func (a *App) Initialize() error {
	if err := a.initDatabase(); err != nil {
		return fmt.Errorf("error inicializando base de datos: %w", err)
	}

	a.initDependencies()

	log.Println("Aplicación inicializada exitosamente")
	return nil
}

func (a *App) initDatabase() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(a.config.DatabaseURI)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return fmt.Errorf("error conectando a MongoDB: %w", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return fmt.Errorf("error verificando conexión a MongoDB: %w", err)
	}

	a.database = client.Database(a.config.DatabaseName)

	log.Printf("Conectado exitosamente a MongoDB: %s", a.config.DatabaseName)
	return nil
}

func (a *App) initDependencies() {

	usuarioRepo := repositories.NewUsuarioRepository(a.database)
	planRepo := repositories.NewPlanRepository(a.database)
	suscripcionRepo := repositories.NewSuscripcionRepository(a.database)

	usuarioService := services.NewUsuarioService(usuarioRepo, a.config)
	planService := services.NewPlanService(planRepo, suscripcionRepo)
	suscripcionService := services.NewSuscripcionService(suscripcionRepo, usuarioRepo, planRepo)

	authMiddleware := middleware.NewAuthMiddleware(a.config.JWTSecret)

	usuarioHandler := v1.NewUsuarioHandler(usuarioService)
	planHandler := v1.NewPlanHandler(planService)
	suscripcionHandler := v1.NewSuscripcionHandler(suscripcionService)

	a.router = v1.NewRouter(
		usuarioHandler,
		planHandler,
		suscripcionHandler,
		authMiddleware,
	)
}

func (a *App) GetRouter() *v1.Router {
	return a.router
}

func (a *App) GetDatabase() *mongo.Database {
	return a.database
}

func (a *App) Close() error {
	if a.database != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := a.database.Client().Disconnect(ctx); err != nil {
			return fmt.Errorf("error cerrando conexión a MongoDB: %w", err)
		}

		log.Println("Conexión a MongoDB cerrada exitosamente")
	}

	return nil
}

func Run(cfg *config.Config) {

	app := NewApp(cfg)

	if err := app.Initialize(); err != nil {
		log.Fatalf("Error inicializando aplicación: %v", err)
	}

	router := app.GetRouter().SetupRoutes()

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.HTTPPort),
		Handler: router,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error iniciando servidor: %v", err)
		}
	}()

	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf(" Error cerrando servidor: %v", err)
	}

	if err := app.Close(); err != nil {
		log.Printf(" Error cerrando aplicación: %v", err)
	}

	log.Println("Servidor cerrado exitosamente")
}
