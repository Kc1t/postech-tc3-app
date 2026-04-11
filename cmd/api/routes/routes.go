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

	v1 := router.Group("/api/v1")

	// --- Auth (publico) ---
	// Protecao contra brute force e feita via account lockout no use case
	// (incrementa failed_attempts, bloqueia conta apos N tentativas).
	// Rate limit por IP foi removido: facilmente burlavel por botnet e
	// deve ser responsabilidade do infra layer (Cloudflare, nginx, WAF).
	authPublic := v1.Group("/auth")
	authPublic.POST("/register", c.AuthHandler.Register)
	authPublic.POST("/login", c.AuthHandler.Login)
	authPublic.POST("/refresh", c.AuthHandler.Refresh)

	// --- Rotas protegidas (qualquer usuario autenticado) ---
	protected := v1.Group("/")
	protected.Use(middleware.Auth(c.Config.JWTSecret))

	// Logout (precisa de JWT)
	protected.POST("/auth/logout", c.AuthHandler.Logout)

	// Consulta de OS — client pode ver (filtra por customerID no handler)
	protected.GET("/service-orders", c.ServiceOrderHandler.FindAll)
	protected.GET("/service-orders/:id", c.ServiceOrderHandler.FindByID)

	// --- Admin only ---
	admin := protected.Group("/")
	admin.Use(middleware.RequireRole("admin"))

	c.CustomerHandler.SetupRoutes(admin)
	c.VehicleHandler.SetupRoutes(admin)
	c.ServiceHandler.SetupRoutes(admin)
	c.PartHandler.SetupRoutes(admin)

	// Service Orders — escrita apenas admin
	orders := admin.Group("/service-orders")
	orders.POST("", c.ServiceOrderHandler.Create)
	orders.PUT("/:id/status", c.ServiceOrderHandler.UpdateStatus)
	orders.PUT("/:id", c.ServiceOrderHandler.Update)
	orders.DELETE("/:id", c.ServiceOrderHandler.Delete)
	admin.GET("/customers/:id/service-orders", c.ServiceOrderHandler.ListByCustomer)
}
