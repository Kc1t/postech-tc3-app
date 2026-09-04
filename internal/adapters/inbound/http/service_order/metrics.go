package serviceorderhandler

import (
	"net/http"

	httputil "github.com/fiap/postech-tc1/internal/adapters/inbound/http"
	"github.com/fiap/postech-tc1/internal/adapters/inbound/http/commands"
	"github.com/gin-gonic/gin"
)

// AverageExecutionTime godoc
// @Summary     Tempo medio de execucao dos servicos
// @Description Media global do intervalo entre aprovacao (started_at) e finalizacao (finished_at) das OSs. Retorna 0 se nenhuma OS foi finalizada ainda.
// @Tags        service-orders
// @Produce     json
// @Success     200 {object} commands.AverageExecutionTimeResponse
// @Failure     500 {object} httputil.ErrorResponse
// @Security    BearerAuth
// @Router      /service-orders/metrics/execution-time [get]
func (h *ServiceOrderHandler) AverageExecutionTime(c *gin.Context) {
	avg, err := h.averageExecutionTime.Execute(c.Request.Context())
	if err != nil {
		httputil.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, commands.AverageExecutionTimeResponse{
		AverageSeconds: avg.Seconds(),
		AverageHuman:   avg.String(),
	})
}
