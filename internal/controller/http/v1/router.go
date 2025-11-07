package v1

import (
	"sw2p2go/internal/middleware"

	_ "sw2p2go/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Router struct {
	usuarioHandler     *UsuarioHandler
	planHandler        *PlanHandler
	suscripcionHandler *SuscripcionHandler
	authMiddleware     *middleware.AuthMiddleware
}

func NewRouter(
	usuarioHandler *UsuarioHandler,
	planHandler *PlanHandler,
	suscripcionHandler *SuscripcionHandler,
	authMiddleware *middleware.AuthMiddleware,
) *Router {
	return &Router{
		usuarioHandler:     usuarioHandler,
		planHandler:        planHandler,
		suscripcionHandler: suscripcionHandler,
		authMiddleware:     authMiddleware,
	}
}

func (r *Router) SetupRoutes() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(r.authMiddleware.CORS())

	// Health Check godoc
	// @Summary      Health Check
	// @Description  Verifica que la API esté funcionando
	// @Tags         System
	// @Accept       json
	// @Produce      json
	// @Success      200  {object}  map[string]string
	// @Router       /health [get]
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "OK", "message": "API is running"})
	})

	v1 := router.Group("/api/v1")

	// Swagger endpoint en /api/v1/docs
	v1.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	///usar el token con cuidau y es /perfil porsi
	auth := v1.Group("/auth")
	{
		auth.POST("/register", r.usuarioHandler.Register)
		auth.POST("/login", r.usuarioHandler.Login)
	}

	userPublic := v1.Group("/usuarios")
	{
		userPublic.GET("", r.usuarioHandler.GetAllUsers)
		userPublic.GET("/:id", r.usuarioHandler.GetUserByID)
	}

	planPublic := v1.Group("/planes")
	{
		planPublic.GET("", r.planHandler.GetAllPlanes)
		planPublic.GET("/:id", r.planHandler.GetPlanByID)
		planPublic.GET("/activos", r.planHandler.GetActivePlanes)
	}

	protected := v1.Group("")
	protected.Use(r.authMiddleware.JWT())
	{
		protectedUsers := protected.Group("/usuarios")
		{
			protectedUsers.PUT("/:id", r.usuarioHandler.UpdateUser)
			protectedUsers.DELETE("/:id", r.usuarioHandler.DeleteUser)
		}

		profile := protected.Group("/perfil")
		{
			profile.GET("", r.usuarioHandler.GetProfile)
		}

		protectedPlans := protected.Group("/planes")
		{
			protectedPlans.POST("", r.planHandler.CreatePlan)
			protectedPlans.PUT("/:id", r.planHandler.UpdatePlan)
			protectedPlans.DELETE("/:id", r.planHandler.DeletePlan)
		}

		suscripciones := protected.Group("/suscripciones")
		{
			suscripciones.POST("", r.suscripcionHandler.CreateSuscripcion)
			suscripciones.GET("", r.suscripcionHandler.GetAllSuscripciones)
			suscripciones.GET("/detalles", r.suscripcionHandler.GetSuscripcionesWithDetails)
			suscripciones.GET("/:id", r.suscripcionHandler.GetSuscripcionByID)
			suscripciones.PUT("/:id", r.suscripcionHandler.UpdateSuscripcion)
			suscripciones.DELETE("/:id", r.suscripcionHandler.CancelSuscripcion)
			suscripciones.GET("/usuario/:user_id", r.suscripcionHandler.GetSuscripcionesByUser)
		}

		misSuscripciones := protected.Group("/mis-suscripciones")
		{
			misSuscripciones.GET("", r.suscripcionHandler.GetMySuscripciones)
		}
	}

	return router
}
