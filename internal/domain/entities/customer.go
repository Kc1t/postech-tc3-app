package entities

import "time"

type Customer struct {
	id        string
	name      string
	document  Document // CPF ou CNPJ (Value Object)
	email     string
	phone     string
	createdAt time.Time
	updatedAt time.Time
}

// NewCustomer cria um novo cliente validando o documento (CPF/CNPJ).
func NewCustomer(name, document, email, phone string) (*Customer, error) {
	doc, err := NewDocument(document)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	return &Customer{
		name:      name,
		document:  doc,
		email:     email,
		phone:     phone,
		createdAt: now,
		updatedAt: now,
	}, nil
}

// ReconstituteCustomer restaura uma entidade a partir de dados persistidos (uso exclusivo de repositories).
func ReconstituteCustomer(id, name, document, email, phone string, createdAt, updatedAt time.Time) *Customer {
	return &Customer{
		id:        id,
		name:      name,
		document:  ReconstituteDocument(document),
		email:     email,
		phone:     phone,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}

func (c *Customer) ID() string           { return c.id }
func (c *Customer) Name() string         { return c.name }
func (c *Customer) Document() string     { return c.document.Value() }
func (c *Customer) DocumentVO() Document { return c.document }
func (c *Customer) Email() string        { return c.email }
func (c *Customer) Phone() string        { return c.phone }
func (c *Customer) CreatedAt() time.Time { return c.createdAt }
func (c *Customer) UpdatedAt() time.Time { return c.updatedAt }

func (c *Customer) SetID(id string)       { c.id = id }
func (c *Customer) SetName(name string)   { c.name = name; c.touch() }
func (c *Customer) SetEmail(email string) { c.email = email; c.touch() }
func (c *Customer) SetPhone(phone string) { c.phone = phone; c.touch() }

func (c *Customer) touch() { c.updatedAt = time.Now() }
