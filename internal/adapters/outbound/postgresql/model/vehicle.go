package pgmodel

import (
	"time"

	"github.com/fiap/postech-tc1/internal/domain/vehicle"
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

func FromVehicle(v *vehicle.Vehicle) *Vehicle {
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

func (m *Vehicle) ToDomain() *vehicle.Vehicle {
	return vehicle.Reconstitute(m.ID, m.CustomerID, m.Plate, m.Brand, m.Model, m.Year, m.CreatedAt, m.UpdatedAt)
}
