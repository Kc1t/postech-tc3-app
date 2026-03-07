package handler

import (
	"net/http"

	"github.com/fiap/postech-tc1/internal/ports"
	"github.com/gin-gonic/gin"
)

type CustomerHandler struct {
	useCase ports.CustomerUseCase
}

func NewCustomerHandler(uc ports.CustomerUseCase) *CustomerHandler {
	return &CustomerHandler{useCase: uc}
}

// Create godoc
// @Summary     Criar cliente
// @Tags        customers
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Router      /customers [post]
func (h *CustomerHandler) Create(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented"})
}

// FindAll godoc
// @Summary     Listar clientes
// @Tags        customers
// @Produce     json
// @Security    BearerAuth
// @Router      /customers [get]
func (h *CustomerHandler) FindAll(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented"})
}

// FindByID godoc
// @Summary     Buscar cliente por ID
// @Tags        customers
// @Produce     json
// @Param       id path string true "Customer ID"
// @Security    BearerAuth
// @Router      /customers/{id} [get]
func (h *CustomerHandler) FindByID(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented"})
}

// Update godoc
// @Summary     Atualizar cliente
// @Tags        customers
// @Accept      json
// @Produce     json
// @Param       id path string true "Customer ID"
// @Security    BearerAuth
// @Router      /customers/{id} [put]
func (h *CustomerHandler) Update(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented"})
}

// Delete godoc
// @Summary     Deletar cliente
// @Tags        customers
// @Produce     json
// @Param       id path string true "Customer ID"
// @Security    BearerAuth
// @Router      /customers/{id} [delete]
func (h *CustomerHandler) Delete(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented"})
}
