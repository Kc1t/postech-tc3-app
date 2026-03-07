package mongomodel

import (
	"time"

	"github.com/fiap/postech-tc1/internal/domain/part"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Part struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `bson:"name"`
	Description string             `bson:"description"`
	Unit        string             `bson:"unit"`
	Price       float64            `bson:"price"`
	Stock       int                `bson:"stock"`
	CreatedAt   time.Time          `bson:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at"`
}

func FromPart(p *part.Part) *Part {
	oid := primitive.NewObjectID()
	if p.ID() != "" {
		oid, _ = primitive.ObjectIDFromHex(p.ID())
	}
	return &Part{
		ID:          oid,
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
	return part.Reconstitute(
		m.ID.Hex(),
		m.Name,
		m.Description,
		m.Unit,
		m.Price,
		m.Stock,
		m.CreatedAt,
		m.UpdatedAt,
	)
}
