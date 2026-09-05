package pgmodel

import (
	"time"

	"github.com/fiap/postech-tc1/internal/domain/entities"
)

type Service struct {
	ID          string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Code        int    `gorm:"not null;uniqueIndex"`
	Name        string `gorm:"not null"`
	Description string
	Price       float64 `gorm:"not null"`
	DurationMin int     `gorm:"not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func FromService(s *entities.Service) *Service {
	return &Service{
		ID:          s.ID(),
		Code:        s.Code(),
		Name:        s.Name(),
		Description: s.Description(),
		Price:       s.Price(),
		DurationMin: s.DurationMin(),
		CreatedAt:   s.CreatedAt(),
		UpdatedAt:   s.UpdatedAt(),
	}
}

func (m *Service) ToDomain() *entities.Service {
	return entities.ReconstituteService(m.ID, m.Code, m.Name, m.Description, m.Price, m.DurationMin, m.CreatedAt, m.UpdatedAt)
}
