package commands

import (
	"time"

	"github.com/fiap/postech-tc1/internal/domain/entities"
)

type CreateServiceOrderRequest struct {
	RequesterDocument string `json:"requester_document" binding:"required"`
	VehiclePlate      string `json:"vehicle_plate"     binding:"required"`
	Notes             string `json:"notes"`
}

type UpdateServiceOrderRequest struct {
	Notes *string `json:"notes"`
}

type ServiceItemRequest struct {
	Code int `json:"code" binding:"required,gt=0"`
}

type PartItemRequest struct {
	ManufacturerCode string `json:"manufacturer_code" binding:"required"`
	Quantity         int    `json:"quantity"           binding:"required,gt=0"`
}

type UpdateStatusRequest struct {
	Status   entities.OrderStatus `json:"status"   binding:"required"`
	Services []ServiceItemRequest `json:"services"`
	Parts    []PartItemRequest    `json:"parts"`
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

type UpdateStatusByCodeRequest struct {
	RequesterDocument string               `json:"requester_document" binding:"required"`
	Status            entities.OrderStatus `json:"status"            binding:"required"`
}

type ServiceOrderResponse struct {
	ID          string                `json:"id"`
	Code        int                   `json:"code"`
	RequesterID string                `json:"requester_id"`
	VehicleID   string                `json:"vehicle_id"`
	Status      entities.OrderStatus  `json:"status"`
	Services    []ServiceItemResponse `json:"services"`
	Parts       []PartItemResponse    `json:"parts"`
	TotalAmount float64               `json:"total_amount"`
	Notes       string                `json:"notes"`
	CreatedAt   string                `json:"created_at"`
	UpdatedAt   string                `json:"updated_at"`
	StartedAt   *string               `json:"started_at,omitempty"`
	FinishedAt  *string               `json:"finished_at,omitempty"`
}

type AverageExecutionTimeResponse struct {
	AverageSeconds float64 `json:"average_seconds" example:"5400"`
	AverageHuman   string  `json:"average_human"   example:"1h30m0s"`
}

func ToServiceOrderResponse(so *entities.ServiceOrder) ServiceOrderResponse {
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

	var startedAt, finishedAt *string
	if t := so.StartedAt(); t != nil {
		s := t.Format(time.RFC3339)
		startedAt = &s
	}
	if t := so.FinishedAt(); t != nil {
		s := t.Format(time.RFC3339)
		finishedAt = &s
	}

	return ServiceOrderResponse{
		ID:          so.ID(),
		Code:        so.Code(),
		RequesterID: so.RequesterID(),
		VehicleID:   so.VehicleID(),
		Status:      so.Status(),
		Services:    services,
		Parts:       parts,
		TotalAmount: so.TotalAmount(),
		Notes:       so.Notes(),
		CreatedAt:   so.CreatedAt().Format(time.RFC3339),
		UpdatedAt:   so.UpdatedAt().Format(time.RFC3339),
		StartedAt:   startedAt,
		FinishedAt:  finishedAt,
	}
}

func ToServiceOrderListResponse(orders []*entities.ServiceOrder) []ServiceOrderResponse {
	resp := make([]ServiceOrderResponse, 0, len(orders))
	for _, so := range orders {
		resp = append(resp, ToServiceOrderResponse(so))
	}
	return resp
}
