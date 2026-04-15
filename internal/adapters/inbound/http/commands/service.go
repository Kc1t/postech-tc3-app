package commands

import (
	"time"

	"github.com/fiap/postech-tc1/internal/domain/entities"
)

type CreateServiceRequest struct {
	Code        int     `json:"code"         binding:"required,gt=0"`
	Name        string  `json:"name"         binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price"        binding:"required,gt=0"`
	DurationMin int     `json:"duration_min" binding:"required,gt=0"`
}

type UpdateServiceRequest struct {
	Code        int     `json:"code"         binding:"omitempty,gt=0"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"        binding:"omitempty,gt=0"`
	DurationMin int     `json:"duration_min" binding:"omitempty,gt=0"`
}

type ServiceResponse struct {
	ID          string  `json:"id"`
	Code        int     `json:"code"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	DurationMin int     `json:"duration_min"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

func (r *CreateServiceRequest) ToDomain() *entities.Service {
	return entities.NewService(r.Code, r.Name, r.Description, r.Price, r.DurationMin)
}

func ToServiceResponse(s *entities.Service) ServiceResponse {
	return ServiceResponse{
		ID:          s.ID(),
		Code:        s.Code(),
		Name:        s.Name(),
		Description: s.Description(),
		Price:       s.Price(),
		DurationMin: s.DurationMin(),
		CreatedAt:   s.CreatedAt().Format(time.RFC3339),
		UpdatedAt:   s.UpdatedAt().Format(time.RFC3339),
	}
}

func ToServiceListResponse(services []*entities.Service) []ServiceResponse {
	resp := make([]ServiceResponse, 0, len(services))
	for _, s := range services {
		resp = append(resp, ToServiceResponse(s))
	}
	return resp
}
