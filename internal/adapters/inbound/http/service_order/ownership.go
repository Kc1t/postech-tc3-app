package serviceorderhandler

import (
	httputil "github.com/fiap/postech-tc1/internal/adapters/inbound/http"
	"github.com/fiap/postech-tc1/internal/domain/entities"
	domainerrors "github.com/fiap/postech-tc1/internal/domain/errors"
	"github.com/gin-gonic/gin"
)

// Token sem CPF so vale para a operacao (admin). O cliente precisa do token emitido por CPF,
// e esse CPF tem que ser o documento pedido.
func authorizeDocument(c *gin.Context, document string) bool {
	authenticated := c.GetString("user_document")
	if authenticated == "" && c.GetString("user_role") == string(entities.RoleAdmin) {
		return true
	}
	if authenticated != "" && authenticated == document {
		return true
	}

	httputil.HandleError(c, domainerrors.ErrDocumentMismatch)
	return false
}
