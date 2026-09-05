package main

import (
	"os"

	"github.com/fiap/postech-tc1/cmd/api/bootstrap"
	"github.com/fiap/postech-tc1/cmd/api/middleware"
	"github.com/fiap/postech-tc1/cmd/api/routes"
	"github.com/fiap/postech-tc1/config"
	"github.com/fiap/postech-tc1/pkg/logger"
	"github.com/gin-gonic/gin"
)

// @title           Workshop API
// @version         1.0
// @description     Sistema Integrado de Atendimento e Execucao de Servicos — Oficina Mecanica
// @host            localhost:8080
// @BasePath        /api/v1
// @securityDefinitions.apikey BearerAuth
// @in              header
// @name            Authorization
func main() {
	container := bootstrap.GetContainer()
	appLogger := logger.New(container.Config.LogLevel)

	if container.Config.AppEnv == config.EnvProduction {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.CorrelationID())
	router.Use(middleware.RequestLogger(appLogger, "/health"))

	routes.Setup(router, container)

	addr := ":" + container.Config.AppPort
	appLogger.Info("server listening", "addr", addr, "env", string(container.Config.AppEnv))

	if err := router.Run(addr); err != nil {
		appLogger.Error("server error", "error", err)
		os.Exit(1)
	}
}
