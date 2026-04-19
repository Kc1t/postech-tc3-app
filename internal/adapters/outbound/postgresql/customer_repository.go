package postgresql

import (
	"context"

	pgmodel "github.com/fiap/postech-tc1/internal/adapters/outbound/postgresql/model"
	"github.com/fiap/postech-tc1/internal/domain/entities"
	"github.com/fiap/postech-tc1/internal/ports"
	"gorm.io/gorm"
)

type customerRepository struct {
	db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) ports.CustomerRepository {
	return &customerRepository{db: db}
}

func (r *customerRepository) Create(ctx context.Context, c *entities.Customer) error {
	m := pgmodel.FromCustomer(c)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return mapError(err)
	}
	c.SetID(m.ID)
	return nil
}

func (r *customerRepository) FindByID(ctx context.Context, id string) (*entities.Customer, error) {
	var m pgmodel.Customer
	if err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		return nil, mapError(err)
	}
	return m.ToDomain(), nil
}

func (r *customerRepository) FindByDocument(ctx context.Context, document string) (*entities.Customer, error) {
	var m pgmodel.Customer
	if err := r.db.WithContext(ctx).First(&m, "document = ?", document).Error; err != nil {
		return nil, mapError(err)
	}
	return m.ToDomain(), nil
}

func (r *customerRepository) FindAll(ctx context.Context) ([]*entities.Customer, error) {
	var docs []pgmodel.Customer
	if err := r.db.WithContext(ctx).Find(&docs).Error; err != nil {
		return nil, mapError(err)
	}
	customers := make([]*entities.Customer, 0, len(docs))
	for i := range docs {
		customers = append(customers, docs[i].ToDomain())
	}
	return customers, nil
}

func (r *customerRepository) Update(ctx context.Context, c *entities.Customer) error {
	m := pgmodel.FromCustomer(c)
	return mapError(r.db.WithContext(ctx).Save(m).Error)
}

func (r *customerRepository) Delete(ctx context.Context, id string) error {
	return mapError(r.db.WithContext(ctx).Delete(&pgmodel.Customer{}, "id = ?", id).Error)
}
