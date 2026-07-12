package main

import (
	"log"

	"github.com/fiap/postech-tc1/cmd/api/bootstrap"
	"github.com/fiap/postech-tc1/cmd/api/routes"
	"github.com/fiap/postech-tc1/config"
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
	// Comentario de teste para validar o fluxo de deploy via CI/CD.
	// Comentario de teste.
	container := bootstrap.GetContainer()

	if container.Config.AppEnv == config.EnvProduction {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	routes.Setup(router, container)

	addr := ":" + container.Config.AppPort
	log.Printf("server listening on %s (env=%s)", addr, container.Config.AppEnv)

	if err := router.Run(addr); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
