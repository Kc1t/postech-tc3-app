package pgmodel

import (
	"time"

	"github.com/fiap/postech-tc1/internal/domain/vehicle"
)

type Vehicle struct {
	ID         string    `db:"id"`
	CustomerID string    `db:"customer_id"`
	Plate      string    `db:"plate"`
	Brand      string    `db:"brand"`
	Model      string    `db:"model"`
	Year       int       `db:"year"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
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
