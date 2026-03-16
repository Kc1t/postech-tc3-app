package pgmodel

import (
	"time"

	"github.com/fiap/postech-tc1/internal/domain/customer"
)

type Customer struct {
	ID        string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name      string    `gorm:"not null"`
	Document  string    `gorm:"uniqueIndex;not null"`
	Email     string    `gorm:"not null"`
	Phone     string    `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func FromCustomer(c *customer.Customer) *Customer {
	return &Customer{
		ID:        c.ID(),
		Name:      c.Name(),
		Document:  c.Document(),
		Email:     c.Email(),
		Phone:     c.Phone(),
		CreatedAt: c.CreatedAt(),
		UpdatedAt: c.UpdatedAt(),
	}
}

func (m *Customer) ToDomain() *customer.Customer {
	return customer.Reconstitute(m.ID, m.Name, m.Document, m.Email, m.Phone, m.CreatedAt, m.UpdatedAt)
}
