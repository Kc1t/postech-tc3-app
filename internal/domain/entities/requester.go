package entities

import "time"

type RequesterStatus string

const (
	RequesterStatusActive   RequesterStatus = "active"
	RequesterStatusInactive RequesterStatus = "inactive"
)

func (s RequesterStatus) IsValid() bool {
	return s == RequesterStatusActive || s == RequesterStatusInactive
}

type Requester struct {
	id        string
	name      string
	document  Document // CPF ou CNPJ (Value Object)
	email     string
	phone     string
	status    RequesterStatus
	createdAt time.Time
	updatedAt time.Time
}

// NewRequester cria um novo cliente validando o documento (CPF/CNPJ).
func NewRequester(name, document, email, phone string) (*Requester, error) {
	doc, err := NewDocument(document)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	return &Requester{
		name:      name,
		document:  doc,
		email:     email,
		phone:     phone,
		status:    RequesterStatusActive,
		createdAt: now,
		updatedAt: now,
	}, nil
}

// ReconstituteRequester restaura uma entidade a partir de dados persistidos (uso exclusivo de repositories).
func ReconstituteRequester(id, name, document, email, phone string, createdAt, updatedAt time.Time) *Requester {
	return &Requester{
		id:        id,
		name:      name,
		document:  ReconstituteDocument(document),
		email:     email,
		phone:     phone,
		status:    RequesterStatusActive,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}

func (c *Requester) ID() string              { return c.id }
func (c *Requester) Name() string            { return c.name }
func (c *Requester) Document() string        { return c.document.Value() }
func (c *Requester) DocumentVO() Document    { return c.document }
func (c *Requester) Email() string           { return c.email }
func (c *Requester) Phone() string           { return c.phone }
func (c *Requester) Status() RequesterStatus { return c.status }
func (c *Requester) IsActive() bool          { return c.status == RequesterStatusActive }
func (c *Requester) CreatedAt() time.Time    { return c.createdAt }
func (c *Requester) UpdatedAt() time.Time    { return c.updatedAt }

func (c *Requester) SetID(id string)       { c.id = id }
func (c *Requester) SetName(name string)   { c.name = name; c.touch() }
func (c *Requester) SetEmail(email string) { c.email = email; c.touch() }
func (c *Requester) SetPhone(phone string) { c.phone = phone; c.touch() }

func (c *Requester) SetStatus(status RequesterStatus) {
	if status.IsValid() {
		c.status = status
	}
}

func (c *Requester) touch() { c.updatedAt = time.Now() }
