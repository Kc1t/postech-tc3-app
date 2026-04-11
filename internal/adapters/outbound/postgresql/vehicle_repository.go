package postgresql

import (
	"context"

	pgmodel "github.com/fiap/postech-tc1/internal/adapters/outbound/postgresql/model"
	"github.com/fiap/postech-tc1/internal/domain/entities"
	"github.com/fiap/postech-tc1/internal/ports"
	"gorm.io/gorm"
)

type vehicleRepository struct {
	db *gorm.DB
}

func NewVehicleRepository(db *gorm.DB) ports.VehicleRepository {
	return &vehicleRepository{db: db}
}

func (r *vehicleRepository) Create(ctx context.Context, v *entities.Vehicle) error {
	m := pgmodel.FromVehicle(v)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	v.SetID(m.ID)
	return nil
}

func (r *vehicleRepository) FindByID(ctx context.Context, id string) (*entities.Vehicle, error) {
	var m pgmodel.Vehicle
	if err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return m.ToDomain(), nil
}

func (r *vehicleRepository) FindByCustomerID(ctx context.Context, customerID string) ([]*entities.Vehicle, error) {
	var docs []pgmodel.Vehicle
	if err := r.db.WithContext(ctx).Find(&docs, "customer_id = ?", customerID).Error; err != nil {
		return nil, err
	}
	vehicles := make([]*entities.Vehicle, 0, len(docs))
	for i := range docs {
		vehicles = append(vehicles, docs[i].ToDomain())
	}
	return vehicles, nil
}

func (r *vehicleRepository) FindAll(ctx context.Context) ([]*entities.Vehicle, error) {
	var docs []pgmodel.Vehicle
	if err := r.db.WithContext(ctx).Find(&docs).Error; err != nil {
		return nil, err
	}
	vehicles := make([]*entities.Vehicle, 0, len(docs))
	for i := range docs {
		vehicles = append(vehicles, docs[i].ToDomain())
	}
	return vehicles, nil
}

func (r *vehicleRepository) Update(ctx context.Context, v *entities.Vehicle) error {
	m := pgmodel.FromVehicle(v)
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *vehicleRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&pgmodel.Vehicle{}, "id = ?", id).Error
}
