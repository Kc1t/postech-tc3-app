package handler

import (
	"net/http"

	"github.com/fiap/postech-tc1/internal/ports"
	"github.com/gin-gonic/gin"
)

type ServiceHandler struct {
	useCase ports.ServiceUseCase
}

func NewServiceHandler(uc ports.ServiceUseCase) *ServiceHandler {
	return &ServiceHandler{useCase: uc}
}

// Create godoc
// @Summary     Criar servico
// @Tags        services
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Router      /services [post]
func (h *ServiceHandler) Create(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented"})
}

// FindAll godoc
// @Summary     Listar servicos
// @Tags        services
// @Produce     json
// @Security    BearerAuth
// @Router      /services [get]
func (h *ServiceHandler) FindAll(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented"})
}

// FindByID godoc
// @Summary     Buscar servico por ID
// @Tags        services
// @Produce     json
// @Param       id path string true "Service ID"
// @Security    BearerAuth
// @Router      /services/{id} [get]
func (h *ServiceHandler) FindByID(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented"})
}

// Update godoc
// @Summary     Atualizar servico
// @Tags        services
// @Accept      json
// @Produce     json
// @Param       id path string true "Service ID"
// @Security    BearerAuth
// @Router      /services/{id} [put]
func (h *ServiceHandler) Update(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented"})
}

// Delete godoc
// @Summary     Deletar servico
// @Tags        services
// @Produce     json
// @Param       id path string true "Service ID"
// @Security    BearerAuth
// @Router      /services/{id} [delete]
func (h *ServiceHandler) Delete(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented"})
}
