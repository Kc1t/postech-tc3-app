package routes

import (
	"net/http"

	"github.com/fiap/postech-tc1/cmd/api/bootstrap"
	"github.com/fiap/postech-tc1/cmd/api/middleware"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/fiap/postech-tc1/docs" // swagger docs gerados pelo swag
)

func Setup(router *gin.Engine, c *bootstrap.Container) {
	router.Use(middleware.CORS())

	// Swagger UI
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Health check (publico)
	router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	prefix := router.Group("/api/v1")

	// Rotas publicas do cliente (sem JWT)
	c.ServiceOrderHandler.SetupPublicRoutes(prefix)

	protected := prefix.Group("/")
	protected.Use(middleware.Auth(c.Config.JWTSecret))

	c.CustomerHandler.SetupRoutes(protected)
	c.VehicleHandler.SetupRoutes(protected)
	c.ServiceOrderHandler.SetupRoutes(protected)
	c.ServiceHandler.SetupRoutes(protected)
	c.PartHandler.SetupRoutes(protected)
}
