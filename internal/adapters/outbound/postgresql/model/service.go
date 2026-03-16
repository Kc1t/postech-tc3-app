package pgmodel

import (
	"time"

	"github.com/fiap/postech-tc1/internal/domain/service"
)

type Service struct {
	ID          string    `db:"id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	Price       float64   `db:"price"`
	DurationMin int       `db:"duration_min"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

func FromService(s *service.Service) *Service {
	return &Service{
		ID:          s.ID(),
		Name:        s.Name(),
		Description: s.Description(),
		Price:       s.Price(),
		DurationMin: s.DurationMin(),
		CreatedAt:   s.CreatedAt(),
		UpdatedAt:   s.UpdatedAt(),
	}
}

func (m *Service) ToDomain() *service.Service {
	return service.Reconstitute(m.ID, m.Name, m.Description, m.Price, m.DurationMin, m.CreatedAt, m.UpdatedAt)
}
