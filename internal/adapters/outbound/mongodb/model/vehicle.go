package mongomodel

import (
	"time"

	"github.com/fiap/postech-tc1/internal/domain/vehicle"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Vehicle struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	CustomerID primitive.ObjectID `bson:"customer_id"`
	Plate      string             `bson:"plate"`
	Brand      string             `bson:"brand"`
	Model      string             `bson:"model"`
	Year       int                `bson:"year"`
	CreatedAt  time.Time          `bson:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at"`
}

func FromVehicle(v *vehicle.Vehicle) *Vehicle {
	oid := primitive.NewObjectID()
	if v.ID() != "" {
		oid, _ = primitive.ObjectIDFromHex(v.ID())
	}
	customerOID, _ := primitive.ObjectIDFromHex(v.CustomerID())
	return &Vehicle{
		ID:         oid,
		CustomerID: customerOID,
		Plate:      v.Plate(),
		Brand:      v.Brand(),
		Model:      v.Model(),
		Year:       v.Year(),
		CreatedAt:  v.CreatedAt(),
		UpdatedAt:  v.UpdatedAt(),
	}
}

func (m *Vehicle) ToDomain() *vehicle.Vehicle {
	return vehicle.Reconstitute(
		m.ID.Hex(),
		m.CustomerID.Hex(),
		m.Plate,
		m.Brand,
		m.Model,
		m.Year,
		m.CreatedAt,
		m.UpdatedAt,
	)
}
