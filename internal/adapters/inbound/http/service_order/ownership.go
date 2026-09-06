package serviceorderhandler

import (
	httputil "github.com/fiap/postech-tc1/internal/adapters/inbound/http"
	domainerrors "github.com/fiap/postech-tc1/internal/domain/errors"
	"github.com/gin-gonic/gin"
)

func authorizeDocument(c *gin.Context, document string) bool {
	authenticated := c.GetString("user_document")
	if authenticated == "" || authenticated == document {
		return true
	}

	httputil.HandleError(c, domainerrors.ErrDocumentMismatch)
	return false
}
