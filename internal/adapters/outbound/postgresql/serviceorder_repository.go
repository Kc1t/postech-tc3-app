package postgresql

import (
	"context"
	"time"

	pgmodel "github.com/fiap/postech-tc1/internal/adapters/outbound/postgresql/model"
	"github.com/fiap/postech-tc1/internal/domain/entities"
	"github.com/fiap/postech-tc1/internal/ports"
	"gorm.io/gorm"
)

type serviceOrderRepository struct {
	db *gorm.DB
}

func NewServiceOrderRepository(db *gorm.DB) ports.ServiceOrderRepository {
	return &serviceOrderRepository{db: db}
}

func (r *serviceOrderRepository) Create(ctx context.Context, so *entities.ServiceOrder) error {
	m := pgmodel.FromServiceOrder(so)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return mapError(err)
	}
	so.SetID(m.ID)
	so.SetCode(m.Code)
	return nil
}

func (r *serviceOrderRepository) FindByID(ctx context.Context, id string) (*entities.ServiceOrder, error) {
	var m pgmodel.ServiceOrder
	if err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		return nil, mapError(err)
	}
	return m.ToDomain(), nil
}

func (r *serviceOrderRepository) FindByCode(ctx context.Context, code int) (*entities.ServiceOrder, error) {
	var m pgmodel.ServiceOrder
	if err := r.db.WithContext(ctx).First(&m, "code = ?", code).Error; err != nil {
		return nil, mapError(err)
	}
	return m.ToDomain(), nil
}

func (r *serviceOrderRepository) FindAll(ctx context.Context) ([]*entities.ServiceOrder, error) {
	var docs []pgmodel.ServiceOrder
	if err := r.db.WithContext(ctx).Find(&docs).Error; err != nil {
		return nil, mapError(err)
	}
	orders := make([]*entities.ServiceOrder, 0, len(docs))
	for i := range docs {
		orders = append(orders, docs[i].ToDomain())
	}
	return orders, nil
}

func (r *serviceOrderRepository) FindByCustomerID(ctx context.Context, customerID string) ([]*entities.ServiceOrder, error) {
	var docs []pgmodel.ServiceOrder
	if err := r.db.WithContext(ctx).Find(&docs, "customer_id = ?", customerID).Error; err != nil {
		return nil, mapError(err)
	}
	orders := make([]*entities.ServiceOrder, 0, len(docs))
	for i := range docs {
		orders = append(orders, docs[i].ToDomain())
	}
	return orders, nil
}

func (r *serviceOrderRepository) UpdateStatus(ctx context.Context, id string, status entities.OrderStatus) error {
	return mapError(r.db.WithContext(ctx).Model(&pgmodel.ServiceOrder{}).Where("id = ?", id).Update("status", status).Error)
}

func (r *serviceOrderRepository) Update(ctx context.Context, so *entities.ServiceOrder) error {
	m := pgmodel.FromServiceOrder(so)
	return mapError(r.db.WithContext(ctx).Save(m).Error)
}

func (r *serviceOrderRepository) Delete(ctx context.Context, id string) error {
	return mapError(r.db.WithContext(ctx).Delete(&pgmodel.ServiceOrder{}, "id = ?", id).Error)
}

// AverageExecutionTime retorna a media global do intervalo entre startedAt
// (momento da aprovacao) e finishedAt (momento da finalizacao). Considera
// apenas OSs com ambos os timestamps gravados. COALESCE garante 0 quando
// nao ha amostras (AVG de nenhum row retorna NULL no PostgreSQL).
func (r *serviceOrderRepository) AverageExecutionTime(ctx context.Context) (time.Duration, error) {
	var averageSeconds float64
	err := r.db.WithContext(ctx).
		Model(&pgmodel.ServiceOrder{}).
		Where("started_at IS NOT NULL AND finished_at IS NOT NULL").
		Select("COALESCE(EXTRACT(EPOCH FROM AVG(finished_at - started_at)), 0)").
		Scan(&averageSeconds).Error
	if err != nil {
		return 0, mapError(err)
	}
	return time.Duration(averageSeconds * float64(time.Second)), nil
}
