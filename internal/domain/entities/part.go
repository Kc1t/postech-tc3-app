package entities

import "time"

type Part struct {
	id               string
	manufacturerCode string // codigo do fabricante impresso na embalagem
	name             string
	description      string
	unit             string // ex: unidade, litro, kg
	price            float64
	stock            int
	createdAt        time.Time
	updatedAt        time.Time
}

func NewPart(manufacturerCode, name, description, unit string, price float64, stock int) *Part {
	now := time.Now()
	return &Part{
		manufacturerCode: manufacturerCode,
		name:             name,
		description:      description,
		unit:             unit,
		price:            price,
		stock:            stock,
		createdAt:        now,
		updatedAt:        now,
	}
}

// ReconstitutePart restaura uma entidade a partir de dados persistidos (uso exclusivo de repositories).
func ReconstitutePart(id, manufacturerCode, name, description, unit string, price float64, stock int, createdAt, updatedAt time.Time) *Part {
	return &Part{
		id:               id,
		manufacturerCode: manufacturerCode,
		name:             name,
		description:      description,
		unit:             unit,
		price:            price,
		stock:            stock,
		createdAt:        createdAt,
		updatedAt:        updatedAt,
	}
}

func (p *Part) ID() string               { return p.id }
func (p *Part) ManufacturerCode() string  { return p.manufacturerCode }
func (p *Part) Name() string              { return p.name }
func (p *Part) Description() string       { return p.description }
func (p *Part) Unit() string              { return p.unit }
func (p *Part) Price() float64            { return p.price }
func (p *Part) Stock() int                { return p.stock }
func (p *Part) CreatedAt() time.Time      { return p.createdAt }
func (p *Part) UpdatedAt() time.Time      { return p.updatedAt }

func (p *Part) SetID(id string)                  { p.id = id }
func (p *Part) SetManufacturerCode(code string)  { p.manufacturerCode = code; p.touch() }
func (p *Part) SetName(n string)                 { p.name = n; p.touch() }
func (p *Part) SetDescription(d string)          { p.description = d; p.touch() }
func (p *Part) SetUnit(u string)                 { p.unit = u; p.touch() }
func (p *Part) SetPrice(pr float64)              { p.price = pr; p.touch() }
func (p *Part) SetStock(s int)                   { p.stock = s; p.touch() }

func (p *Part) touch() { p.updatedAt = time.Now() }
