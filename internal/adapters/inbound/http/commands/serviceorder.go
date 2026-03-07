package commands

import (
	"time"

	"github.com/fiap/postech-tc1/internal/domain/serviceorder"
)

type ServiceItemRequest struct {
	ServiceID   string  `json:"service_id"  binding:"required"`
	Description string  `json:"description" binding:"required"`
	Price       float64 `json:"price"       binding:"required,gt=0"`
}

type PartItemRequest struct {
	PartID      string  `json:"part_id"     binding:"required"`
	Description string  `json:"description" binding:"required"`
	Quantity    int     `json:"quantity"    binding:"required,gt=0"`
	UnitPrice   float64 `json:"unit_price"  binding:"required,gt=0"`
}

type CreateServiceOrderRequest struct {
	CustomerID string               `json:"customer_id" binding:"required"`
	VehicleID  string               `json:"vehicle_id"  binding:"required"`
	Services   []ServiceItemRequest `json:"services"`
	Parts      []PartItemRequest    `json:"parts"`
	Notes      string               `json:"notes"`
}

type UpdateStatusRequest struct {
	Status serviceorder.Status `json:"status" binding:"required"`
}

type ServiceItemResponse struct {
	ServiceID   string  `json:"service_id"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

type PartItemResponse struct {
	PartID      string  `json:"part_id"`
	Description string  `json:"description"`
	Quantity    int     `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
}

type ServiceOrderResponse struct {
	ID          string                `json:"id"`
	CustomerID  string                `json:"customer_id"`
	VehicleID   string                `json:"vehicle_id"`
	Status      serviceorder.Status   `json:"status"`
	Services    []ServiceItemResponse `json:"services"`
	Parts       []PartItemResponse    `json:"parts"`
	TotalAmount float64               `json:"total_amount"`
	Notes       string                `json:"notes"`
	CreatedAt   string                `json:"created_at"`
	UpdatedAt   string                `json:"updated_at"`
}

func (r *CreateServiceOrderRequest) ToDomain() *serviceorder.ServiceOrder {
	so := serviceorder.New(r.CustomerID, r.VehicleID)
	so.SetNotes(r.Notes)
	for _, s := range r.Services {
		so.AddService(serviceorder.ServiceItem{
			ServiceID:   s.ServiceID,
			Description: s.Description,
			Price:       s.Price,
		})
	}
	for _, p := range r.Parts {
		so.AddPart(serviceorder.PartItem{
			PartID:      p.PartID,
			Description: p.Description,
			Quantity:    p.Quantity,
			UnitPrice:   p.UnitPrice,
		})
	}
	return so
}

func ToServiceOrderResponse(so *serviceorder.ServiceOrder) ServiceOrderResponse {
	services := make([]ServiceItemResponse, 0, len(so.Services()))
	for _, s := range so.Services() {
		services = append(services, ServiceItemResponse{
			ServiceID:   s.ServiceID,
			Description: s.Description,
			Price:       s.Price,
		})
	}

	parts := make([]PartItemResponse, 0, len(so.Parts()))
	for _, p := range so.Parts() {
		parts = append(parts, PartItemResponse{
			PartID:      p.PartID,
			Description: p.Description,
			Quantity:    p.Quantity,
			UnitPrice:   p.UnitPrice,
		})
	}

	return ServiceOrderResponse{
		ID:          so.ID(),
		CustomerID:  so.CustomerID(),
		VehicleID:   so.VehicleID(),
		Status:      so.Status(),
		Services:    services,
		Parts:       parts,
		TotalAmount: so.TotalAmount(),
		Notes:       so.Notes(),
		CreatedAt:   so.CreatedAt().Format(time.RFC3339),
		UpdatedAt:   so.UpdatedAt().Format(time.RFC3339),
	}
}

func ToServiceOrderListResponse(orders []*serviceorder.ServiceOrder) []ServiceOrderResponse {
	resp := make([]ServiceOrderResponse, 0, len(orders))
	for _, so := range orders {
		resp = append(resp, ToServiceOrderResponse(so))
	}
	return resp
}
