package pgmodel

import (
	"time"

	"github.com/fiap/postech-tc1/internal/domain/entities"
)

type Requester struct {
	ID        string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name      string    `gorm:"not null"`
	Document  string    `gorm:"uniqueIndex;not null"`
	Email     string    `gorm:"not null"`
	Phone     string    `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func FromRequester(c *entities.Requester) *Requester {
	return &Requester{
		ID:        c.ID(),
		Name:      c.Name(),
		Document:  c.Document(),
		Email:     c.Email(),
		Phone:     c.Phone(),
		CreatedAt: c.CreatedAt(),
		UpdatedAt: c.UpdatedAt(),
	}
}

func (m *Requester) ToDomain() *entities.Requester {
	return entities.ReconstituteRequester(m.ID, m.Name, m.Document, m.Email, m.Phone, m.CreatedAt, m.UpdatedAt)
}
