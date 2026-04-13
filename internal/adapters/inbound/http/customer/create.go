package customerhandler

import (
	"errors"
	"net/http"

	"github.com/fiap/postech-tc1/internal/adapters/inbound/http/commands"
	domainerrors "github.com/fiap/postech-tc1/internal/domain/errors"
	"github.com/gin-gonic/gin"
)

// Create godoc
// @Summary     Criar cliente
// @Tags        customers
// @Accept      json
// @Produce     json
// @Param       body body commands.CreateCustomerRequest true "Dados do cliente"
// @Success     201 {object} commands.CustomerResponse
// @Security    BearerAuth
// @Router      /customers [post]
func (h *CustomerHandler) Create(c *gin.Context) {
	var req commands.CreateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	customer := req.ToDomain()
	if err := h.create.Execute(c.Request.Context(), customer); err != nil {
		if errors.Is(err, domainerrors.ErrAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"error": "customer already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, commands.ToCustomerResponse(customer))
}
