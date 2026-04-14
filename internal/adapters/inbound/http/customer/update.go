package customerhandler

import (
	"net/http"

	httputil "github.com/fiap/postech-tc1/internal/adapters/inbound/http"
	"github.com/fiap/postech-tc1/internal/adapters/inbound/http/commands"
	"github.com/gin-gonic/gin"
)

// Update godoc
// @Summary     Atualizar cliente
// @Tags        customers
// @Accept      json
// @Produce     json
// @Param       id path string true "Customer ID"
// @Param       body body commands.UpdateCustomerRequest true "Dados para atualizar"
// @Success     200 {object} commands.CustomerResponse
// @Security    BearerAuth
// @Router      /customers/{id} [put]
func (h *CustomerHandler) Update(c *gin.Context) {
	var req commands.UpdateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.HandleBadRequest(c, err)
		return
	}

	customer, err := h.getByID.Execute(c.Request.Context(), c.Param("id"))
	if err != nil {
		httputil.HandleError(c, err, msgNotFound)
		return
	}

	if req.Name != "" {
		customer.SetName(req.Name)
	}
	if req.Email != "" {
		customer.SetEmail(req.Email)
	}
	if req.Phone != "" {
		customer.SetPhone(req.Phone)
	}

	if err := h.update.Execute(c.Request.Context(), customer); err != nil {
		httputil.HandleError(c, err)
		return
	}


	c.JSON(http.StatusOK, commands.ToCustomerResponse(customer))
}
