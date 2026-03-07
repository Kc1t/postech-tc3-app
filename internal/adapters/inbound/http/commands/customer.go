package commands

import (
	"time"

	"github.com/fiap/postech-tc1/internal/domain/customer"
)

type CreateCustomerRequest struct {
	Name     string `json:"name"     binding:"required"`
	Document string `json:"document" binding:"required"`
	Email    string `json:"email"    binding:"required,email"`
	Phone    string `json:"phone"`
}

type UpdateCustomerRequest struct {
	Name  string `json:"name"`
	Email string `json:"email" binding:"omitempty,email"`
	Phone string `json:"phone"`
}

type CustomerResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Document  string `json:"document"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func (r *CreateCustomerRequest) ToDomain() *customer.Customer {
	return customer.New(r.Name, r.Document, r.Email, r.Phone)
}

func ToCustomerResponse(c *customer.Customer) CustomerResponse {
	return CustomerResponse{
		ID:        c.ID(),
		Name:      c.Name(),
		Document:  c.Document(),
		Email:     c.Email(),
		Phone:     c.Phone(),
		CreatedAt: c.CreatedAt().Format(time.RFC3339),
		UpdatedAt: c.UpdatedAt().Format(time.RFC3339),
	}
}

func ToCustomerListResponse(customers []*customer.Customer) []CustomerResponse {
	resp := make([]CustomerResponse, 0, len(customers))
	for _, c := range customers {
		resp = append(resp, ToCustomerResponse(c))
	}
	return resp
}
