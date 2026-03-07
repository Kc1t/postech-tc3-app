package handler

import (
	"net/http"

	"github.com/fiap/postech-tc1/internal/ports"
	"github.com/gin-gonic/gin"
)

type PartHandler struct {
	useCase ports.PartUseCase
}

func NewPartHandler(uc ports.PartUseCase) *PartHandler {
	return &PartHandler{useCase: uc}
}

// Create godoc
// @Summary     Cadastrar peca/insumo
// @Tags        parts
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Router      /parts [post]
func (h *PartHandler) Create(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented"})
}

// FindAll godoc
// @Summary     Listar pecas/insumos
// @Tags        parts
// @Produce     json
// @Security    BearerAuth
// @Router      /parts [get]
func (h *PartHandler) FindAll(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented"})
}

// FindByID godoc
// @Summary     Buscar peca/insumo por ID
// @Tags        parts
// @Produce     json
// @Param       id path string true "Part ID"
// @Security    BearerAuth
// @Router      /parts/{id} [get]
func (h *PartHandler) FindByID(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented"})
}

// Update godoc
// @Summary     Atualizar peca/insumo
// @Tags        parts
// @Accept      json
// @Produce     json
// @Param       id path string true "Part ID"
// @Security    BearerAuth
// @Router      /parts/{id} [put]
func (h *PartHandler) Update(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented"})
}

// Delete godoc
// @Summary     Deletar peca/insumo
// @Tags        parts
// @Produce     json
// @Param       id path string true "Part ID"
// @Security    BearerAuth
// @Router      /parts/{id} [delete]
func (h *PartHandler) Delete(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented"})
}
