package commands

import (
	"time"

	"github.com/fiap/postech-tc1/internal/domain/entities"
)

type CreateVehicleRequest struct {
	CustomerDocument string `json:"customer_document" binding:"required"`
	Plate            string `json:"plate"             binding:"required"`
	Brand            string `json:"brand"             binding:"required"`
	Model            string `json:"model"             binding:"required"`
	Year             int    `json:"year"              binding:"required"`
}

type UpdateVehicleRequest struct {
	Plate string `json:"plate"`
	Brand string `json:"brand"`
	Model string `json:"model"`
	Year  int    `json:"year"`
}

type VehicleResponse struct {
	ID         string `json:"id"`
	CustomerID string `json:"customer_id"`
	Plate      string `json:"plate"`
	Brand      string `json:"brand"`
	Model      string `json:"model"`
	Year       int    `json:"year"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

func (r *CreateVehicleRequest) ToDomain() *entities.Vehicle {
	return entities.NewVehicle("", r.Plate, r.Brand, r.Model, r.Year)
}

func ToVehicleResponse(v *entities.Vehicle) VehicleResponse {
	return VehicleResponse{
		ID:         v.ID(),
		CustomerID: v.CustomerID(),
		Plate:      v.Plate(),
		Brand:      v.Brand(),
		Model:      v.Model(),
		Year:       v.Year(),
		CreatedAt:  v.CreatedAt().Format(time.RFC3339),
		UpdatedAt:  v.UpdatedAt().Format(time.RFC3339),
	}
}

func ToVehicleListResponse(vehicles []*entities.Vehicle) []VehicleResponse {
	resp := make([]VehicleResponse, 0, len(vehicles))
	for _, v := range vehicles {
		resp = append(resp, ToVehicleResponse(v))
	}
	return resp
}
