package pgmodel

import (
	"time"

	"github.com/fiap/postech-tc1/internal/domain/entities"
)

type Vehicle struct {
	ID         string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	CustomerID string    `gorm:"type:uuid;not null;index"`
	Plate      string    `gorm:"uniqueIndex;not null"`
	Brand      string    `gorm:"not null"`
	Model      string    `gorm:"not null"`
	Year       int       `gorm:"not null"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func FromVehicle(v *entities.Vehicle) *Vehicle {
	return &Vehicle{
		ID:         v.ID(),
		CustomerID: v.CustomerID(),
		Plate:      v.Plate(),
		Brand:      v.Brand(),
		Model:      v.Model(),
		Year:       v.Year(),
		CreatedAt:  v.CreatedAt(),
		UpdatedAt:  v.UpdatedAt(),
	}
}

func (m *Vehicle) ToDomain() *entities.Vehicle {
	return entities.ReconstituteVehicle(m.ID, m.CustomerID, m.Plate, m.Brand, m.Model, m.Year, m.CreatedAt, m.UpdatedAt)
}
