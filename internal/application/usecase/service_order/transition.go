package serviceorderuc

import (
	"time"

	"github.com/fiap/postech-tc1/internal/domain/entities"
)

// statusTransition guarda de qual status a OS saiu e desde quando estava nele. O updated_at da OS
// marca a ultima alteracao, que no fluxo normal e a entrada no status atual.
type statusTransition struct {
	from  entities.OrderStatus
	since time.Time
}

func captureTransition(so *entities.ServiceOrder) statusTransition {
	return statusTransition{from: so.Status(), since: so.UpdatedAt()}
}

func (t statusTransition) logAttrs(now time.Time) []any {
	attrs := []any{"from_status", string(t.from)}
	if !t.since.IsZero() && now.After(t.since) {
		attrs = append(attrs, "seconds_in_status", now.Sub(t.since).Seconds())
	}
	return attrs
}
