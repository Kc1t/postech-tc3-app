package authhandler

import (
	"net/http"

	"github.com/fiap/postech-tc1/internal/adapters/inbound/http/commands"
	"github.com/gin-gonic/gin"
)

// Logout godoc
// @Summary     Encerrar sessao (revoga todos os refresh tokens)
// @Tags        auth
// @Produce     json
// @Security    BearerAuth
// @Success     200 {object} commands.MessageResponse
// @Failure     401 {object} commands.MessageResponse
// @Router      /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not identified"})
		return
	}

	userID, ok := userIDVal.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user context"})
		return
	}

	if err := h.logout.Execute(c.Request.Context(), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, commands.MessageResponse{Message: "logged out successfully"})
}
