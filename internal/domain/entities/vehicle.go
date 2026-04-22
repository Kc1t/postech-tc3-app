package entities

import "time"

// VehicleInput carrega os dados necessarios para cadastrar um veiculo.
type VehicleInput struct {
	CustomerDocument string
	Plate            string
	Brand            string
	Model            string
	Year             int
}

type Vehicle struct {
	id         string
	customerID string
	plate      Plate // Value Object
	brand      string
	model      string
	year       int
	createdAt  time.Time
	updatedAt  time.Time
}

// NewVehicle cria um novo veiculo validando a placa (formato antigo ou Mercosul).
func NewVehicle(customerID, plate, brand, model string, year int) (*Vehicle, error) {
	p, err := NewPlate(plate)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	return &Vehicle{
		customerID: customerID,
		plate:      p,
		brand:      brand,
		model:      model,
		year:       year,
		createdAt:  now,
		updatedAt:  now,
	}, nil
}

// ReconstituteVehicle restaura uma entidade a partir de dados persistidos (uso exclusivo de repositories).
func ReconstituteVehicle(id, customerID, plate, brand, model string, year int, createdAt, updatedAt time.Time) *Vehicle {
	return &Vehicle{
		id:         id,
		customerID: customerID,
		plate:      ReconstitutePlate(plate),
		brand:      brand,
		model:      model,
		year:       year,
		createdAt:  createdAt,
		updatedAt:  updatedAt,
	}
}

func (v *Vehicle) ID() string           { return v.id }
func (v *Vehicle) CustomerID() string   { return v.customerID }
func (v *Vehicle) Plate() string        { return v.plate.Value() }
func (v *Vehicle) PlateVO() Plate       { return v.plate }
func (v *Vehicle) Brand() string        { return v.brand }
func (v *Vehicle) Model() string        { return v.model }
func (v *Vehicle) Year() int            { return v.year }
func (v *Vehicle) CreatedAt() time.Time { return v.createdAt }
func (v *Vehicle) UpdatedAt() time.Time { return v.updatedAt }

func (v *Vehicle) SetID(id string)         { v.id = id }
func (v *Vehicle) SetCustomerID(id string) { v.customerID = id }

// SetPlate valida e atualiza a placa.
func (v *Vehicle) SetPlate(raw string) error {
	p, err := NewPlate(raw)
	if err != nil {
		return err
	}
	v.plate = p
	v.touch()
	return nil
}
func (v *Vehicle) SetBrand(b string) { v.brand = b; v.touch() }
func (v *Vehicle) SetModel(m string) { v.model = m; v.touch() }
func (v *Vehicle) SetYear(y int)     { v.year = y; v.touch() }

func (v *Vehicle) touch() { v.updatedAt = time.Now() }
