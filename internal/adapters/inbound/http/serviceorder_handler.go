package handler

import (
	"net/http"

	"github.com/fiap/postech-tc1/internal/ports"
	"github.com/gin-gonic/gin"
)

type ServiceOrderHandler struct {
	useCase ports.ServiceOrderUseCase
}

func NewServiceOrderHandler(uc ports.ServiceOrderUseCase) *ServiceOrderHandler {
	return &ServiceOrderHandler{useCase: uc}
}

// Create godoc
// @Summary     Criar ordem de servico
// @Tags        service-orders
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Router      /service-orders [post]
func (h *ServiceOrderHandler) Create(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented"})
}

// FindAll godoc
// @Summary     Listar ordens de servico
// @Tags        service-orders
// @Produce     json
// @Security    BearerAuth
// @Router      /service-orders [get]
func (h *ServiceOrderHandler) FindAll(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented"})
}

// FindByID godoc
// @Summary     Buscar ordem de servico por ID
// @Tags        service-orders
// @Produce     json
// @Param       id path string true "ServiceOrder ID"
// @Security    BearerAuth
// @Router      /service-orders/{id} [get]
func (h *ServiceOrderHandler) FindByID(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented"})
}

// UpdateStatus godoc
// @Summary     Atualizar status da ordem de servico
// @Tags        service-orders
// @Accept      json
// @Produce     json
// @Param       id path string true "ServiceOrder ID"
// @Security    BearerAuth
// @Router      /service-orders/{id}/status [put]
func (h *ServiceOrderHandler) UpdateStatus(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented"})
}

// Delete godoc
// @Summary     Deletar ordem de servico
// @Tags        service-orders
// @Produce     json
// @Param       id path string true "ServiceOrder ID"
// @Security    BearerAuth
// @Router      /service-orders/{id} [delete]
func (h *ServiceOrderHandler) Delete(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented"})
}
