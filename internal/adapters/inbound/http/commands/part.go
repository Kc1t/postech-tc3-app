package commands

import (
	"time"

	"github.com/fiap/postech-tc1/internal/domain/entities"
)

type CreatePartRequest struct {
	Name        string  `json:"name"        binding:"required"`
	Description string  `json:"description"`
	Unit        string  `json:"unit"        binding:"required"`
	Price       float64 `json:"price"       binding:"required,gt=0"`
	Stock       int     `json:"stock"       binding:"min=0"`
}

type UpdatePartRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Unit        string  `json:"unit"`
	Price       float64 `json:"price"  binding:"omitempty,gt=0"`
	Stock       int     `json:"stock"  binding:"omitempty,min=0"`
}

type AdjustStockRequest struct {
	Delta int `json:"delta" binding:"required"` // positivo = entrada, negativo = saida
}

type PartResponse struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Unit        string  `json:"unit"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

func (r *CreatePartRequest) ToDomain() *entities.Part {
	return entities.NewPart(r.Name, r.Description, r.Unit, r.Price, r.Stock)
}

func ToPartResponse(p *entities.Part) PartResponse {
	return PartResponse{
		ID:          p.ID(),
		Name:        p.Name(),
		Description: p.Description(),
		Unit:        p.Unit(),
		Price:       p.Price(),
		Stock:       p.Stock(),
		CreatedAt:   p.CreatedAt().Format(time.RFC3339),
		UpdatedAt:   p.UpdatedAt().Format(time.RFC3339),
	}
}

func ToPartListResponse(parts []*entities.Part) []PartResponse {
	resp := make([]PartResponse, 0, len(parts))
	for _, p := range parts {
		resp = append(resp, ToPartResponse(p))
	}
	return resp
}
