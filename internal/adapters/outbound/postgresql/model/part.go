package pgmodel

import (
	"time"

	"github.com/fiap/postech-tc1/internal/domain/entities"
)

type Part struct {
	ID               string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	ManufacturerCode string `gorm:"not null;uniqueIndex"`
	Name             string `gorm:"not null"`
	Description      string
	Unit             string  `gorm:"not null"`
	Price            float64 `gorm:"not null"`
	Stock            int     `gorm:"not null;default:0"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

func FromPart(p *entities.Part) *Part {
	return &Part{
		ID:               p.ID(),
		ManufacturerCode: p.ManufacturerCode(),
		Name:             p.Name(),
		Description:      p.Description(),
		Unit:             p.Unit(),
		Price:            p.Price(),
		Stock:            p.Stock(),
		CreatedAt:        p.CreatedAt(),
		UpdatedAt:        p.UpdatedAt(),
	}
}

func (m *Part) ToDomain() *entities.Part {
	return entities.ReconstitutePart(m.ID, m.ManufacturerCode, m.Name, m.Description, m.Unit, m.Price, m.Stock, m.CreatedAt, m.UpdatedAt)
}
