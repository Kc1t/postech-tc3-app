package postgresql

import (
	"context"

	pgmodel "github.com/fiap/postech-tc1/internal/adapters/outbound/postgresql/model"
	"github.com/fiap/postech-tc1/internal/domain/entities"
	"github.com/fiap/postech-tc1/internal/ports"
	"gorm.io/gorm"
)

type requesterRepository struct {
	db *gorm.DB
}

func NewRequesterRepository(db *gorm.DB) ports.RequesterRepository {
	return &requesterRepository{db: db}
}

func (r *requesterRepository) Create(ctx context.Context, c *entities.Requester) error {
	m := pgmodel.FromRequester(c)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return mapError(err)
	}
	c.SetID(m.ID)
	return nil
}

func (r *requesterRepository) FindByID(ctx context.Context, id string) (*entities.Requester, error) {
	var m pgmodel.Requester
	if err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		return nil, mapError(err)
	}
	return m.ToDomain(), nil
}

func (r *requesterRepository) FindByDocument(ctx context.Context, document string) (*entities.Requester, error) {
	var m pgmodel.Requester
	if err := r.db.WithContext(ctx).First(&m, "document = ?", document).Error; err != nil {
		return nil, mapError(err)
	}
	return m.ToDomain(), nil
}

func (r *requesterRepository) FindAll(ctx context.Context) ([]*entities.Requester, error) {
	var docs []pgmodel.Requester
	if err := r.db.WithContext(ctx).Find(&docs).Error; err != nil {
		return nil, mapError(err)
	}
	requesters := make([]*entities.Requester, 0, len(docs))
	for i := range docs {
		requesters = append(requesters, docs[i].ToDomain())
	}
	return requesters, nil
}

func (r *requesterRepository) Update(ctx context.Context, c *entities.Requester) error {
	m := pgmodel.FromRequester(c)
	return mapError(r.db.WithContext(ctx).Save(m).Error)
}

func (r *requesterRepository) Delete(ctx context.Context, id string) error {
	return mapError(r.db.WithContext(ctx).Delete(&pgmodel.Requester{}, "id = ?", id).Error)
}
