package entities

import (
	"time"

	domainerrors "github.com/fiap/postech-tc1/internal/domain/errors"
)

type OrderStatus string

const (
	StatusReceived         OrderStatus = "received"
	StatusInDiagnosis      OrderStatus = "in_diagnosis"
	StatusAwaitingApproval OrderStatus = "awaiting_approval"
	StatusInExecution      OrderStatus = "in_execution"
	StatusFinished         OrderStatus = "finished"
	StatusDelivered        OrderStatus = "delivered"
)

// validTransitions define a maquina de estados da Ordem de Servico.
var validTransitions = map[OrderStatus][]OrderStatus{
	StatusReceived:         {StatusInDiagnosis},
	StatusInDiagnosis:      {StatusAwaitingApproval},
	StatusAwaitingApproval: {StatusInExecution, StatusReceived},
	StatusInExecution:      {StatusFinished},
	StatusFinished:         {StatusDelivered},
	StatusDelivered:        {},
}

// allStatuses permite validar se um status informado e conhecido.
var allStatuses = map[OrderStatus]bool{
	StatusReceived:         true,
	StatusInDiagnosis:      true,
	StatusAwaitingApproval: true,
	StatusInExecution:      true,
	StatusFinished:         true,
	StatusDelivered:        true,
}

// IsValidStatus verifica se o status informado e um valor conhecido.
func IsValidStatus(s OrderStatus) bool {
	return allStatuses[s]
}

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

func (so *ServiceOrder) SetID(id string)   { so.id = id }
func (so *ServiceOrder) SetNotes(n string) { so.notes = n; so.touch() }

// UpdateStatus valida a transicao de status antes de aplicar.
func (so *ServiceOrder) UpdateStatus(s OrderStatus) error {
	if !IsValidStatus(s) {
		return domainerrors.ErrInvalidStatusValue
	}
	allowed, exists := validTransitions[so.status]
	if !exists {
		return domainerrors.ErrInvalidStatus
	}
	for _, a := range allowed {
		if a == s {
			so.status = s
			so.touch()
			return nil
		}
	}
	return domainerrors.ErrInvalidStatus
}

// SetServices substitui a lista de servicos e recalcula o total.
func (so *ServiceOrder) SetServices(items []ServiceItem) {
	so.services = items
	so.recalcTotal()
}

// SetParts substitui a lista de pecas e recalcula o total.
func (so *ServiceOrder) SetParts(items []PartItem) {
	so.parts = items
	so.recalcTotal()
}

func (so *ServiceOrder) AddService(item ServiceItem) {
	so.services = append(so.services, item)
	so.recalcTotal()
}

func (so *ServiceOrder) AddPart(item PartItem) {
	so.parts = append(so.parts, item)
	so.recalcTotal()
}

// RecalcTotal forca o recalculo do total (uso apos ajustes externos).
func (so *ServiceOrder) RecalcTotal() {
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
