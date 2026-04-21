package pgmodel

import (
	"time"

	"github.com/fiap/postech-tc1/internal/domain/entities"
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
	ID          string              `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Code        int                 `gorm:"not null;uniqueIndex;autoIncrement"`
	CustomerID  string              `gorm:"type:uuid;not null;index;constraint:OnDelete:RESTRICT"`
	Customer    Customer            `gorm:"foreignKey:CustomerID"`
	VehicleID   string              `gorm:"type:uuid;not null;constraint:OnDelete:RESTRICT"`
	Vehicle     Vehicle             `gorm:"foreignKey:VehicleID"`
	Status      entities.OrderStatus `gorm:"not null"`
	Services    []ServiceItem       `gorm:"serializer:json"`
	Parts       []PartItem          `gorm:"serializer:json"`
	TotalAmount float64             `gorm:"not null"`
	Notes       string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func FromServiceOrder(so *entities.ServiceOrder) *ServiceOrder {
	services := make([]ServiceItem, 0, len(so.Services()))
	for _, s := range so.Services() {
		services = append(services, ServiceItem{ServiceID: s.ServiceID, Description: s.Description, Price: s.Price})
	}

	parts := make([]PartItem, 0, len(so.Parts()))
	for _, p := range so.Parts() {
		parts = append(parts, PartItem{PartID: p.PartID, Description: p.Description, Quantity: p.Quantity, UnitPrice: p.UnitPrice})
	}

	return &ServiceOrder{
		ID:          so.ID(),
		Code:        so.Code(),
		CustomerID:  so.CustomerID(),
		VehicleID:   so.VehicleID(),
		Status:      so.Status(),
		Services:    services,
		Parts:       parts,
		TotalAmount: so.TotalAmount(),
		Notes:       so.Notes(),
		CreatedAt:   so.CreatedAt(),
		UpdatedAt:   so.UpdatedAt(),
	}
}

func (m *ServiceOrder) ToDomain() *entities.ServiceOrder {
	services := make([]entities.ServiceItem, 0, len(m.Services))
	for _, s := range m.Services {
		services = append(services, entities.ServiceItem{ServiceID: s.ServiceID, Description: s.Description, Price: s.Price})
	}

	parts := make([]entities.PartItem, 0, len(m.Parts))
	for _, p := range m.Parts {
		parts = append(parts, entities.PartItem{PartID: p.PartID, Description: p.Description, Quantity: p.Quantity, UnitPrice: p.UnitPrice})
	}

	return entities.ReconstituteServiceOrder(
		m.ID, m.Code, m.CustomerID, m.VehicleID, m.Status,
		services, parts, m.TotalAmount, m.Notes, m.CreatedAt, m.UpdatedAt,
	)
}
