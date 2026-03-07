package mongomodel

import (
	"time"

	"github.com/fiap/postech-tc1/internal/domain/serviceorder"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ServiceItem struct {
	ServiceID   primitive.ObjectID `bson:"service_id"`
	Description string             `bson:"description"`
	Price       float64            `bson:"price"`
}

type PartItem struct {
	PartID      primitive.ObjectID `bson:"part_id"`
	Description string             `bson:"description"`
	Quantity    int                `bson:"quantity"`
	UnitPrice   float64            `bson:"unit_price"`
}

type ServiceOrder struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	CustomerID  primitive.ObjectID `bson:"customer_id"`
	VehicleID   primitive.ObjectID `bson:"vehicle_id"`
	Status      serviceorder.Status `bson:"status"`
	Services    []ServiceItem      `bson:"services"`
	Parts       []PartItem         `bson:"parts"`
	TotalAmount float64            `bson:"total_amount"`
	Notes       string             `bson:"notes"`
	CreatedAt   time.Time          `bson:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at"`
}

func FromServiceOrder(so *serviceorder.ServiceOrder) *ServiceOrder {
	oid := primitive.NewObjectID()
	if so.ID() != "" {
		oid, _ = primitive.ObjectIDFromHex(so.ID())
	}
	customerOID, _ := primitive.ObjectIDFromHex(so.CustomerID())
	vehicleOID, _ := primitive.ObjectIDFromHex(so.VehicleID())

	services := make([]ServiceItem, 0, len(so.Services()))
	for _, s := range so.Services() {
		sOID, _ := primitive.ObjectIDFromHex(s.ServiceID)
		services = append(services, ServiceItem{
			ServiceID:   sOID,
			Description: s.Description,
			Price:       s.Price,
		})
	}

	parts := make([]PartItem, 0, len(so.Parts()))
	for _, p := range so.Parts() {
		pOID, _ := primitive.ObjectIDFromHex(p.PartID)
		parts = append(parts, PartItem{
			PartID:      pOID,
			Description: p.Description,
			Quantity:    p.Quantity,
			UnitPrice:   p.UnitPrice,
		})
	}

	return &ServiceOrder{
		ID:          oid,
		CustomerID:  customerOID,
		VehicleID:   vehicleOID,
		Status:      so.Status(),
		Services:    services,
		Parts:       parts,
		TotalAmount: so.TotalAmount(),
		Notes:       so.Notes(),
		CreatedAt:   so.CreatedAt(),
		UpdatedAt:   so.UpdatedAt(),
	}
}

func (m *ServiceOrder) ToDomain() *serviceorder.ServiceOrder {
	services := make([]serviceorder.ServiceItem, 0, len(m.Services))
	for _, s := range m.Services {
		services = append(services, serviceorder.ServiceItem{
			ServiceID:   s.ServiceID.Hex(),
			Description: s.Description,
			Price:       s.Price,
		})
	}

	parts := make([]serviceorder.PartItem, 0, len(m.Parts))
	for _, p := range m.Parts {
		parts = append(parts, serviceorder.PartItem{
			PartID:      p.PartID.Hex(),
			Description: p.Description,
			Quantity:    p.Quantity,
			UnitPrice:   p.UnitPrice,
		})
	}

	return serviceorder.Reconstitute(
		m.ID.Hex(),
		m.CustomerID.Hex(),
		m.VehicleID.Hex(),
		m.Status,
		services,
		parts,
		m.TotalAmount,
		m.Notes,
		m.CreatedAt,
		m.UpdatedAt,
	)
}
