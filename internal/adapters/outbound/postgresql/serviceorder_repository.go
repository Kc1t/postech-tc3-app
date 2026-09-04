package postgresql

import (
	"context"
	"time"

	pgmodel "github.com/fiap/postech-tc1/internal/adapters/outbound/postgresql/model"
	"github.com/fiap/postech-tc1/internal/domain/entities"
	domainerrors "github.com/fiap/postech-tc1/internal/domain/errors"
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
	err := r.db.WithContext(ctx).
		Where("status NOT IN ?", []string{
			string(entities.StatusFinished),
			string(entities.StatusDelivered),
		}).
		Order("CASE status" +
			" WHEN 'in_execution' THEN 1" +
			" WHEN 'awaiting_approval' THEN 2" +
			" WHEN 'in_diagnosis' THEN 3" +
			" WHEN 'received' THEN 4" +
			" END").
		Order("created_at ASC").
		Find(&docs).Error
	if err != nil {
		return nil, mapError(err)
	}
	orders := make([]*entities.ServiceOrder, 0, len(docs))
	for i := range docs {
		orders = append(orders, docs[i].ToDomain())
	}
	return orders, nil
}

func (r *serviceOrderRepository) FindByRequesterID(ctx context.Context, requesterID string) ([]*entities.ServiceOrder, error) {
	var docs []pgmodel.ServiceOrder
	if err := r.db.WithContext(ctx).Find(&docs, "requester_id = ?", requesterID).Error; err != nil {
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

func (r *serviceOrderRepository) ApplyApprovalTransition(ctx context.Context, so *entities.ServiceOrder) error {
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, p := range so.Parts() {
			if err := decrementStockTx(tx, p.PartID, p.Quantity); err != nil {
				return err
			}
		}
		m := pgmodel.FromServiceOrder(so)
		return tx.Save(m).Error
	})
	return mapError(err)
}

func decrementStockTx(tx *gorm.DB, partID string, quantity int) error {
	result := tx.Model(&pgmodel.Part{}).
		Where("id = ? AND stock - ? >= 0", partID, quantity).
		UpdateColumn("stock", gorm.Expr("stock - ?", quantity))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		var count int64
		if err := tx.Model(&pgmodel.Part{}).Where("id = ?", partID).Count(&count).Error; err != nil {
			return err
		}
		if count == 0 {
			return domainerrors.ErrNotFound
		}
		return domainerrors.ErrInsufficientStock
	}
	return nil
}

func (r *serviceOrderRepository) Delete(ctx context.Context, id string) error {
	return mapError(r.db.WithContext(ctx).Delete(&pgmodel.ServiceOrder{}, "id = ?", id).Error)
}

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
