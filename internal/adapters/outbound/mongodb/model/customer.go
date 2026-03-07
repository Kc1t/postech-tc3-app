package mongomodel

import (
	"time"

	"github.com/fiap/postech-tc1/internal/domain/customer"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Customer struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name"`
	Document  string             `bson:"document"`
	Email     string             `bson:"email"`
	Phone     string             `bson:"phone"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

func FromCustomer(c *customer.Customer) *Customer {
	oid := primitive.NewObjectID()
	if c.ID() != "" {
		oid, _ = primitive.ObjectIDFromHex(c.ID())
	}
	return &Customer{
		ID:        oid,
		Name:      c.Name(),
		Document:  c.Document(),
		Email:     c.Email(),
		Phone:     c.Phone(),
		CreatedAt: c.CreatedAt(),
		UpdatedAt: c.UpdatedAt(),
	}
}

func (m *Customer) ToDomain() *customer.Customer {
	return customer.Reconstitute(
		m.ID.Hex(),
		m.Name,
		m.Document,
		m.Email,
		m.Phone,
		m.CreatedAt,
		m.UpdatedAt,
	)
}
