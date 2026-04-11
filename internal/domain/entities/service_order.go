package entities

import "time"

type OrderStatus string

const (
	StatusReceived         OrderStatus = "received"
	StatusInDiagnosis      OrderStatus = "in_diagnosis"
	StatusAwaitingApproval OrderStatus = "awaiting_approval"
	StatusInExecution      OrderStatus = "in_execution"
	StatusFinished         OrderStatus = "finished"
	StatusDelivered        OrderStatus = "delivered"
)

// ServiceItem e PartItem sao value objects — identificados por valor, sem identidade propria.
type ServiceItem struct {
	ServiceID   string
	Description string
	Price       float64
}

type PartItem struct {
	PartID      string
	Description string
	Quantity    int
	UnitPrice   float64
}

type ServiceOrder struct {
	id          string
	customerID  string
	vehicleID   string
	status      OrderStatus
	services    []ServiceItem
	parts       []PartItem
	totalAmount float64
	notes       string
	createdAt   time.Time
	updatedAt   time.Time
}

func NewServiceOrder(customerID, vehicleID string) *ServiceOrder {
	now := time.Now()
	return &ServiceOrder{
		customerID: customerID,
		vehicleID:  vehicleID,
		status:     StatusReceived,
		services:   []ServiceItem{},
		parts:      []PartItem{},
		createdAt:  now,
		updatedAt:  now,
	}
}

// ReconstituteServiceOrder restaura uma entidade a partir de dados persistidos (uso exclusivo de repositories).
func ReconstituteServiceOrder(
	id, customerID, vehicleID string,
	status OrderStatus,
	services []ServiceItem,
	parts []PartItem,
	totalAmount float64,
	notes string,
	createdAt, updatedAt time.Time,
) *ServiceOrder {
	return &ServiceOrder{
		id:          id,
		customerID:  customerID,
		vehicleID:   vehicleID,
		status:      status,
		services:    services,
		parts:       parts,
		totalAmount: totalAmount,
		notes:       notes,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
	}
}

func (so *ServiceOrder) ID() string              { return so.id }
func (so *ServiceOrder) CustomerID() string      { return so.customerID }
func (so *ServiceOrder) VehicleID() string       { return so.vehicleID }
func (so *ServiceOrder) Status() OrderStatus     { return so.status }
func (so *ServiceOrder) Services() []ServiceItem { return so.services }
func (so *ServiceOrder) Parts() []PartItem       { return so.parts }
func (so *ServiceOrder) TotalAmount() float64    { return so.totalAmount }
func (so *ServiceOrder) Notes() string           { return so.notes }
func (so *ServiceOrder) CreatedAt() time.Time    { return so.createdAt }
func (so *ServiceOrder) UpdatedAt() time.Time    { return so.updatedAt }

func (so *ServiceOrder) SetID(id string)           { so.id = id }
func (so *ServiceOrder) SetNotes(n string)         { so.notes = n; so.touch() }
func (so *ServiceOrder) UpdateStatus(s OrderStatus) { so.status = s; so.touch() }

func (so *ServiceOrder) AddService(item ServiceItem) {
	so.services = append(so.services, item)
	so.recalcTotal()
}

func (so *ServiceOrder) AddPart(item PartItem) {
	so.parts = append(so.parts, item)
	so.recalcTotal()
}

func (so *ServiceOrder) recalcTotal() {
	total := 0.0
	for _, s := range so.services {
		total += s.Price
	}
	for _, p := range so.parts {
		total += float64(p.Quantity) * p.UnitPrice
	}
	so.totalAmount = total
	so.touch()
}

func (so *ServiceOrder) touch() { so.updatedAt = time.Now() }
