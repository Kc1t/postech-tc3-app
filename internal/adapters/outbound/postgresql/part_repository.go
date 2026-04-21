package postgresql

import (
	"context"

	pgmodel "github.com/fiap/postech-tc1/internal/adapters/outbound/postgresql/model"
	"github.com/fiap/postech-tc1/internal/domain/entities"
	"github.com/fiap/postech-tc1/internal/ports"
	"gorm.io/gorm"
)

type partRepository struct {
	db *gorm.DB
}

func NewPartRepository(db *gorm.DB) ports.PartRepository {
	return &partRepository{db: db}
}

func (r *partRepository) Create(ctx context.Context, p *entities.Part) error {
	m := pgmodel.FromPart(p)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return mapError(err)
	}
	p.SetID(m.ID)
	return nil
}

func (r *partRepository) FindByID(ctx context.Context, id string) (*entities.Part, error) {
	var m pgmodel.Part
	if err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		return nil, mapError(err)
	}
	return m.ToDomain(), nil
}

func (r *partRepository) FindByIDs(ctx context.Context, ids []string) ([]*entities.Part, error) {
	var docs []pgmodel.Part
	if err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&docs).Error; err != nil {
		return nil, mapError(err)
	}
	parts := make([]*entities.Part, 0, len(docs))
	for i := range docs {
		parts = append(parts, docs[i].ToDomain())
	}
	return parts, nil
}

func (r *partRepository) FindByManufacturerCodes(ctx context.Context, codes []string) ([]*entities.Part, error) {
	var docs []pgmodel.Part
	if err := r.db.WithContext(ctx).Where("manufacturer_code IN ?", codes).Find(&docs).Error; err != nil {
		return nil, mapError(err)
	}
	parts := make([]*entities.Part, 0, len(docs))
	for i := range docs {
		parts = append(parts, docs[i].ToDomain())
	}
	return parts, nil
}

func (r *partRepository) FindAll(ctx context.Context) ([]*entities.Part, error) {
	var docs []pgmodel.Part
	if err := r.db.WithContext(ctx).Find(&docs).Error; err != nil {
		return nil, mapError(err)
	}
	parts := make([]*entities.Part, 0, len(docs))
	for i := range docs {
		parts = append(parts, docs[i].ToDomain())
	}
	return parts, nil
}

func (r *partRepository) Update(ctx context.Context, p *entities.Part) error {
	m := pgmodel.FromPart(p)
	return mapError(r.db.WithContext(ctx).Save(m).Error)
}

func (r *partRepository) Delete(ctx context.Context, id string) error {
	return mapError(r.db.WithContext(ctx).Delete(&pgmodel.Part{}, "id = ?", id).Error)
}

func (r *partRepository) UpdateStock(ctx context.Context, id string, delta int) error {
	return mapError(r.db.WithContext(ctx).Model(&pgmodel.Part{}).Where("id = ?", id).UpdateColumn("stock", gorm.Expr("stock + ?", delta)).Error)
}
