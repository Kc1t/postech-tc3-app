package mongomodel

import (
	"time"

	"github.com/fiap/postech-tc1/internal/domain/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `bson:"name"`
	Description string             `bson:"description"`
	Price       float64            `bson:"price"`
	DurationMin int                `bson:"duration_min"`
	CreatedAt   time.Time          `bson:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at"`
}

func FromService(s *service.Service) *Service {
	oid := primitive.NewObjectID()
	if s.ID() != "" {
		oid, _ = primitive.ObjectIDFromHex(s.ID())
	}
	return &Service{
		ID:          oid,
		Name:        s.Name(),
		Description: s.Description(),
		Price:       s.Price(),
		DurationMin: s.DurationMin(),
		CreatedAt:   s.CreatedAt(),
		UpdatedAt:   s.UpdatedAt(),
	}
}

func (m *Service) ToDomain() *service.Service {
	return service.Reconstitute(
		m.ID.Hex(),
		m.Name,
		m.Description,
		m.Price,
		m.DurationMin,
		m.CreatedAt,
		m.UpdatedAt,
	)
}
