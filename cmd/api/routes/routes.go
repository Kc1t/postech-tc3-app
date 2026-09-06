package routes

import (
	"net/http"

	_ "github.com/fiap/postech-tc1/docs"

	"github.com/fiap/postech-tc1/cmd/api/bootstrap"
	"github.com/fiap/postech-tc1/cmd/api/middleware"
	"github.com/fiap/postech-tc1/internal/domain/entities"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Setup(router *gin.Engine, c *bootstrap.Container) {
	router.Use(middleware.CORS())
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	prefix := router.Group("/api/v1")

	// Rotas publicas
	auth := prefix.Group("/auth")
	auth.POST("/register", c.AuthHandler.Register)
	auth.POST("/login", c.AuthHandler.Login)
	auth.POST("/refresh", c.AuthHandler.Refresh)

	// Rotas autenticadas (qualquer role)
	protected := prefix.Group("/")
	protected.Use(middleware.Auth(c.Config.JWTSecret))
	protected.POST("/auth/logout", c.AuthHandler.Logout)

	// Rotas do cliente: exigem o JWT emitido a partir do CPF pela lambda de
	// autenticacao, e so respondem para o documento dono do token.
	c.ServiceOrderHandler.SetupRequesterRoutes(protected)

	// Rotas protegidas (JWT obrigatorio) e ADMIN
	admin := prefix.Group("/")
	admin.Use(middleware.Auth(c.Config.JWTSecret))
	admin.Use(middleware.RequireRole(string(entities.RoleAdmin)))

	c.RequesterHandler.SetupRoutes(admin)
	c.VehicleHandler.SetupRoutes(admin)
	c.ServiceHandler.SetupRoutes(admin)
	c.PartHandler.SetupRoutes(admin)
	c.ServiceOrderHandler.SetupRoutes(admin)
}
