package postgresql

import (
	"context"

	pgmodel "github.com/fiap/postech-tc1/internal/adapters/outbound/postgresql/model"
	"github.com/fiap/postech-tc1/internal/domain/entities"
	"github.com/fiap/postech-tc1/internal/ports"
	"gorm.io/gorm"
)

type serviceRepository struct {
	db *gorm.DB
}

func NewServiceRepository(db *gorm.DB) ports.ServiceRepository {
	return &serviceRepository{db: db}
}

func (r *serviceRepository) Create(ctx context.Context, s *entities.Service) error {
	m := pgmodel.FromService(s)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return mapError(err)
	}
	s.SetID(m.ID)
	return nil
}

func (r *serviceRepository) FindByID(ctx context.Context, id string) (*entities.Service, error) {
	var m pgmodel.Service
	if err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		return nil, mapError(err)
	}
	return m.ToDomain(), nil
}

func (r *serviceRepository) FindByIDs(ctx context.Context, ids []string) ([]*entities.Service, error) {
	var docs []pgmodel.Service
	if err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&docs).Error; err != nil {
		return nil, mapError(err)
	}
	services := make([]*entities.Service, 0, len(docs))
	for i := range docs {
		services = append(services, docs[i].ToDomain())
	}
	return services, nil
}

func (r *serviceRepository) FindByCodes(ctx context.Context, codes []int) ([]*entities.Service, error) {
	var docs []pgmodel.Service
	if err := r.db.WithContext(ctx).Where("code IN ?", codes).Find(&docs).Error; err != nil {
		return nil, mapError(err)
	}
	services := make([]*entities.Service, 0, len(docs))
	for i := range docs {
		services = append(services, docs[i].ToDomain())
	}
	return services, nil
}

func (r *serviceRepository) FindAll(ctx context.Context) ([]*entities.Service, error) {
	var docs []pgmodel.Service
	if err := r.db.WithContext(ctx).Find(&docs).Error; err != nil {
		return nil, mapError(err)
	}
	services := make([]*entities.Service, 0, len(docs))
	for i := range docs {
		services = append(services, docs[i].ToDomain())
	}
	return services, nil
}

func (r *serviceRepository) Update(ctx context.Context, s *entities.Service) error {
	m := pgmodel.FromService(s)
	return mapError(r.db.WithContext(ctx).Save(m).Error)
}

func (r *serviceRepository) Delete(ctx context.Context, id string) error {
	return mapError(r.db.WithContext(ctx).Delete(&pgmodel.Service{}, "id = ?", id).Error)
}
