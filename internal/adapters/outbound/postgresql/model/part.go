package pgmodel

import (
	"time"

	"github.com/fiap/postech-tc1/internal/domain/part"
)

type Part struct {
	ID          string    `db:"id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	Unit        string    `db:"unit"`
	Price       float64   `db:"price"`
	Stock       int       `db:"stock"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

func FromPart(p *part.Part) *Part {
	return &Part{
		ID:          p.ID(),
		Name:        p.Name(),
		Description: p.Description(),
		Unit:        p.Unit(),
		Price:       p.Price(),
		Stock:       p.Stock(),
		CreatedAt:   p.CreatedAt(),
		UpdatedAt:   p.UpdatedAt(),
	}
}

func (m *Part) ToDomain() *part.Part {
	return part.Reconstitute(m.ID, m.Name, m.Description, m.Unit, m.Price, m.Stock, m.CreatedAt, m.UpdatedAt)
}
