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

	// Rotas protegidas por JWT
	//TODO fazer um setuproutes dentro de Handler pra cada recurso, e aqui so chamar c.CustomerHandler.SetupRoutes(customers) por exemplo
	protected := v1.Group("/")
	protected.Use(middleware.Auth(c.Config.JWTSecret))
	{
		// Clientes
		customers := protected.Group("/customers")
		{
			customers.POST("", c.CustomerHandler.Create)
			customers.GET("", c.CustomerHandler.FindAll)
			customers.GET("/:id", c.CustomerHandler.FindByID)
			customers.PUT("/:id", c.CustomerHandler.Update)
			customers.DELETE("/:id", c.CustomerHandler.Delete)
		}

		// Veiculos
		vehicles := protected.Group("/vehicles")
		{
			vehicles.POST("", c.VehicleHandler.Create)
			vehicles.GET("", c.VehicleHandler.FindAll)
			vehicles.GET("/:id", c.VehicleHandler.FindByID)
			vehicles.PUT("/:id", c.VehicleHandler.Update)
			vehicles.DELETE("/:id", c.VehicleHandler.Delete)
		}

		// Ordens de Servico
		orders := protected.Group("/service-orders")
		{
			orders.POST("", c.ServiceOrderHandler.Create)
			orders.GET("", c.ServiceOrderHandler.FindAll)
			orders.GET("/:id", c.ServiceOrderHandler.FindByID)
			orders.PUT("/:id/status", c.ServiceOrderHandler.UpdateStatus)
			orders.DELETE("/:id", c.ServiceOrderHandler.Delete)
		}

		// Servicos
		services := protected.Group("/services")
		{
			services.POST("", c.ServiceHandler.Create)
			services.GET("", c.ServiceHandler.FindAll)
			services.GET("/:id", c.ServiceHandler.FindByID)
			services.PUT("/:id", c.ServiceHandler.Update)
			services.DELETE("/:id", c.ServiceHandler.Delete)
		}

		// Pecas e Insumos
		parts := protected.Group("/parts")
		{
			parts.POST("", c.PartHandler.Create)
			parts.GET("", c.PartHandler.FindAll)
			parts.GET("/:id", c.PartHandler.FindByID)
			parts.PUT("/:id", c.PartHandler.Update)
			parts.DELETE("/:id", c.PartHandler.Delete)
		}
	}
}
