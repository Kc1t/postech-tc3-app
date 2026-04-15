package entities

import "time"

type Service struct {
	id          string
	code        int // codigo numerico sequencial (ex: 1, 2, 3)
	name        string
	description string
	price       float64
	durationMin int // tempo estimado em minutos
	createdAt   time.Time
	updatedAt   time.Time
}

func NewService(code int, name, description string, price float64, durationMin int) *Service {
	now := time.Now()
	return &Service{
		code:        code,
		name:        name,
		description: description,
		price:       price,
		durationMin: durationMin,
		createdAt:   now,
		updatedAt:   now,
	}
}

// ReconstituteService restaura uma entidade a partir de dados persistidos (uso exclusivo de repositories).
func ReconstituteService(id string, code int, name, description string, price float64, durationMin int, createdAt, updatedAt time.Time) *Service {
	return &Service{
		id:          id,
		code:        code,
		name:        name,
		description: description,
		price:       price,
		durationMin: durationMin,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
	}
}

func (s *Service) ID() string           { return s.id }
func (s *Service) Code() int            { return s.code }
func (s *Service) Name() string         { return s.name }
func (s *Service) Description() string  { return s.description }
func (s *Service) Price() float64       { return s.price }
func (s *Service) DurationMin() int     { return s.durationMin }
func (s *Service) CreatedAt() time.Time { return s.createdAt }
func (s *Service) UpdatedAt() time.Time { return s.updatedAt }

func (s *Service) SetID(id string)         { s.id = id }
func (s *Service) SetCode(c int)           { s.code = c; s.touch() }
func (s *Service) SetName(n string)        { s.name = n; s.touch() }
func (s *Service) SetDescription(d string) { s.description = d; s.touch() }
func (s *Service) SetPrice(p float64)      { s.price = p; s.touch() }
func (s *Service) SetDurationMin(d int)    { s.durationMin = d; s.touch() }

func (s *Service) touch() { s.updatedAt = time.Now() }
