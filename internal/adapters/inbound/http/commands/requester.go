package commands

import (
	"time"

	"github.com/fiap/postech-tc1/internal/domain/entities"
)

type CreateRequesterRequest struct {
	Name     string `json:"name"     binding:"required"`
	Document string `json:"document" binding:"required"`
	Email    string `json:"email"    binding:"required,email"`
	Phone    string `json:"phone"`
}

type UpdateRequesterRequest struct {
	Name  string `json:"name"`
	Email string `json:"email" binding:"omitempty,email"`
	Phone string `json:"phone"`
}

type RequesterResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Document  string `json:"document"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func (r *CreateRequesterRequest) ToDomain() (*entities.Requester, error) {
	return entities.NewRequester(r.Name, r.Document, r.Email, r.Phone)
}

func ToRequesterResponse(c *entities.Requester) RequesterResponse {
	return RequesterResponse{
		ID:        c.ID(),
		Name:      c.Name(),
		Document:  c.Document(),
		Email:     c.Email(),
		Phone:     c.Phone(),
		CreatedAt: c.CreatedAt().Format(time.RFC3339),
		UpdatedAt: c.UpdatedAt().Format(time.RFC3339),
	}
}

func ToRequesterListResponse(requesters []*entities.Requester) []RequesterResponse {
	resp := make([]RequesterResponse, 0, len(requesters))
	for _, c := range requesters {
		resp = append(resp, ToRequesterResponse(c))
	}
	return resp
}
