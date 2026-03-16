package pgmodel

import (
	"encoding/json"
	"time"

	"github.com/fiap/postech-tc1/internal/domain/serviceorder"
)

type ServiceItem struct {
	ServiceID   string  `json:"service_id"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

type PartItem struct {
	PartID      string  `json:"part_id"`
	Description string  `json:"description"`
	Quantity    int     `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
}

type ServiceOrder struct {
	ID          string                `db:"id"`
	CustomerID  string                `db:"customer_id"`
	VehicleID   string                `db:"vehicle_id"`
	Status      serviceorder.Status   `db:"status"`
	Services    []ServiceItem         `db:"-"`
	ServicesRaw []byte                `db:"services"`
	Parts       []PartItem            `db:"-"`
	PartsRaw    []byte                `db:"parts"`
	TotalAmount float64               `db:"total_amount"`
	Notes       string                `db:"notes"`
	CreatedAt   time.Time             `db:"created_at"`
	UpdatedAt   time.Time             `db:"updated_at"`
}

func FromServiceOrder(so *serviceorder.ServiceOrder) *ServiceOrder {
	services := make([]ServiceItem, 0, len(so.Services()))
	for _, s := range so.Services() {
		services = append(services, ServiceItem{
			ServiceID:   s.ServiceID,
			Description: s.Description,
			Price:       s.Price,
		})
	}

	parts := make([]PartItem, 0, len(so.Parts()))
	for _, p := range so.Parts() {
		parts = append(parts, PartItem{
			PartID:      p.PartID,
			Description: p.Description,
			Quantity:    p.Quantity,
			UnitPrice:   p.UnitPrice,
		})
	}

	servicesRaw, _ := json.Marshal(services)
	partsRaw, _ := json.Marshal(parts)

	return &ServiceOrder{
		ID:          so.ID(),
		CustomerID:  so.CustomerID(),
		VehicleID:   so.VehicleID(),
		Status:      so.Status(),
		Services:    services,
		ServicesRaw: servicesRaw,
		Parts:       parts,
		PartsRaw:    partsRaw,
		TotalAmount: so.TotalAmount(),
		Notes:       so.Notes(),
		CreatedAt:   so.CreatedAt(),
		UpdatedAt:   so.UpdatedAt(),
	}
}

func (m *ServiceOrder) ToDomain() *serviceorder.ServiceOrder {
	_ = json.Unmarshal(m.ServicesRaw, &m.Services)
	_ = json.Unmarshal(m.PartsRaw, &m.Parts)

	services := make([]serviceorder.ServiceItem, 0, len(m.Services))
	for _, s := range m.Services {
		services = append(services, serviceorder.ServiceItem{
			ServiceID:   s.ServiceID,
			Description: s.Description,
			Price:       s.Price,
		})
	}

	parts := make([]serviceorder.PartItem, 0, len(m.Parts))
	for _, p := range m.Parts {
		parts = append(parts, serviceorder.PartItem{
			PartID:      p.PartID,
			Description: p.Description,
			Quantity:    p.Quantity,
			UnitPrice:   p.UnitPrice,
		})
	}

	return serviceorder.Reconstitute(
		m.ID, m.CustomerID, m.VehicleID, m.Status,
		services, parts, m.TotalAmount, m.Notes, m.CreatedAt, m.UpdatedAt,
	)
}
