package postgresql

import (
	"context"

	pgmodel "github.com/fiap/postech-tc1/internal/adapters/outbound/postgresql/model"
	"github.com/fiap/postech-tc1/internal/domain/service"
	"github.com/fiap/postech-tc1/internal/ports"
	"github.com/jackc/pgx/v5/pgxpool"
)

type serviceRepository struct {
	db *pgxpool.Pool
}

func NewServiceRepository(db *pgxpool.Pool) ports.ServiceRepository {
	return &serviceRepository{db: db}
}

func (r *serviceRepository) Create(ctx context.Context, s *service.Service) error {
	query := `
		INSERT INTO services (name, description, price, duration_min, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`
	var id string
	err := r.db.QueryRow(ctx, query, s.Name(), s.Description(), s.Price(), s.DurationMin(), s.CreatedAt(), s.UpdatedAt()).Scan(&id)
	if err != nil {
		return err
	}
	s.SetID(id)
	return nil
}

func (r *serviceRepository) FindByID(ctx context.Context, id string) (*service.Service, error) {
	query := `SELECT id, name, description, price, duration_min, created_at, updated_at FROM services WHERE id=$1`
	var m pgmodel.Service
	err := r.db.QueryRow(ctx, query, id).Scan(&m.ID, &m.Name, &m.Description, &m.Price, &m.DurationMin, &m.CreatedAt, &m.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return m.ToDomain(), nil
}

func (r *serviceRepository) FindAll(ctx context.Context) ([]*service.Service, error) {
	query := `SELECT id, name, description, price, duration_min, created_at, updated_at FROM services`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var services []*service.Service
	for rows.Next() {
		var m pgmodel.Service
		if err := rows.Scan(&m.ID, &m.Name, &m.Description, &m.Price, &m.DurationMin, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, err
		}
		services = append(services, m.ToDomain())
	}
	return services, rows.Err()
}

func (r *serviceRepository) Update(ctx context.Context, s *service.Service) error {
	query := `UPDATE services SET name=$1, description=$2, price=$3, duration_min=$4, updated_at=$5 WHERE id=$6`
	_, err := r.db.Exec(ctx, query, s.Name(), s.Description(), s.Price(), s.DurationMin(), s.UpdatedAt(), s.ID())
	return err
}

func (r *serviceRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM services WHERE id=$1`, id)
	return err
}
