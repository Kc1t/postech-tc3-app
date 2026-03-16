package pgmodel

import (
	"time"

	"github.com/fiap/postech-tc1/internal/domain/customer"
)

type Customer struct {
	ID        string    `db:"id"`
	Name      string    `db:"name"`
	Document  string    `db:"document"`
	Email     string    `db:"email"`
	Phone     string    `db:"phone"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
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
