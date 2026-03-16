package pgmodel

import (
	"time"

	"github.com/fiap/postech-tc1/internal/domain/service"
)

type Service struct {
	ID          string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name        string    `gorm:"not null"`
	Description string
	Price       float64   `gorm:"not null"`
	DurationMin int       `gorm:"not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
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
