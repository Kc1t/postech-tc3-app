package customerhandler

import (
	"net/http"

	httputil "github.com/fiap/postech-tc1/internal/adapters/inbound/http"
	"github.com/fiap/postech-tc1/internal/adapters/inbound/http/commands"
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
		httputil.HandleBadRequest(c, err)
		return
	}

	customer := req.ToDomain()
	if err := h.create.Execute(c.Request.Context(), customer); err != nil {
		httputil.HandleError(c, err, msgAlreadyExists)
		return
	}

	c.JSON(http.StatusCreated, commands.ToCustomerResponse(customer))
}
