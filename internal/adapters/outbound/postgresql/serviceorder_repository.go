package postgresql

import (
	"context"

	pgmodel "github.com/fiap/postech-tc1/internal/adapters/outbound/postgresql/model"
	"github.com/fiap/postech-tc1/internal/domain/serviceorder"
	"github.com/fiap/postech-tc1/internal/ports"
	"gorm.io/gorm"
)

type serviceOrderRepository struct {
	db *gorm.DB
}

func NewServiceOrderRepository(db *gorm.DB) ports.ServiceOrderRepository {
	return &serviceOrderRepository{db: db}
}

func (r *serviceOrderRepository) Create(ctx context.Context, so *serviceorder.ServiceOrder) error {
	m := pgmodel.FromServiceOrder(so)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	so.SetID(m.ID)
	return nil
}

func (r *serviceOrderRepository) FindByID(ctx context.Context, id string) (*serviceorder.ServiceOrder, error) {
	var m pgmodel.ServiceOrder
	if err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return m.ToDomain(), nil
}

func (r *serviceOrderRepository) FindAll(ctx context.Context) ([]*serviceorder.ServiceOrder, error) {
	var docs []pgmodel.ServiceOrder
	if err := r.db.WithContext(ctx).Find(&docs).Error; err != nil {
		return nil, err
	}
	orders := make([]*serviceorder.ServiceOrder, 0, len(docs))
	for i := range docs {
		orders = append(orders, docs[i].ToDomain())
	}
	return orders, nil
}

func (r *serviceOrderRepository) FindByCustomerID(ctx context.Context, customerID string) ([]*serviceorder.ServiceOrder, error) {
	var docs []pgmodel.ServiceOrder
	if err := r.db.WithContext(ctx).Find(&docs, "customer_id = ?", customerID).Error; err != nil {
		return nil, err
	}
	orders := make([]*serviceorder.ServiceOrder, 0, len(docs))
	for i := range docs {
		orders = append(orders, docs[i].ToDomain())
	}
	return orders, nil
}

func (r *serviceOrderRepository) UpdateStatus(ctx context.Context, id string, status serviceorder.Status) error {
	return r.db.WithContext(ctx).Model(&pgmodel.ServiceOrder{}).Where("id = ?", id).Update("status", status).Error
}

func (r *serviceOrderRepository) Update(ctx context.Context, so *serviceorder.ServiceOrder) error {
	m := pgmodel.FromServiceOrder(so)
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *serviceOrderRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&pgmodel.ServiceOrder{}, "id = ?", id).Error
}
