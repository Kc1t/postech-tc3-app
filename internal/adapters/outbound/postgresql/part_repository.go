package postgresql

import (
	"context"

	pgmodel "github.com/fiap/postech-tc1/internal/adapters/outbound/postgresql/model"
	"github.com/fiap/postech-tc1/internal/domain/part"
	"github.com/fiap/postech-tc1/internal/ports"
	"github.com/jackc/pgx/v5/pgxpool"
)

type partRepository struct {
	db *pgxpool.Pool
}

func NewPartRepository(db *pgxpool.Pool) ports.PartRepository {
	return &partRepository{db: db}
}

func (r *partRepository) Create(ctx context.Context, p *part.Part) error {
	query := `
		INSERT INTO parts (name, description, unit, price, stock, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id`
	var id string
	err := r.db.QueryRow(ctx, query, p.Name(), p.Description(), p.Unit(), p.Price(), p.Stock(), p.CreatedAt(), p.UpdatedAt()).Scan(&id)
	if err != nil {
		return err
	}
	p.SetID(id)
	return nil
}

func (r *partRepository) FindByID(ctx context.Context, id string) (*part.Part, error) {
	query := `SELECT id, name, description, unit, price, stock, created_at, updated_at FROM parts WHERE id=$1`
	var m pgmodel.Part
	err := r.db.QueryRow(ctx, query, id).Scan(&m.ID, &m.Name, &m.Description, &m.Unit, &m.Price, &m.Stock, &m.CreatedAt, &m.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return m.ToDomain(), nil
}

func (r *partRepository) FindAll(ctx context.Context) ([]*part.Part, error) {
	query := `SELECT id, name, description, unit, price, stock, created_at, updated_at FROM parts`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var parts []*part.Part
	for rows.Next() {
		var m pgmodel.Part
		if err := rows.Scan(&m.ID, &m.Name, &m.Description, &m.Unit, &m.Price, &m.Stock, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, err
		}
		parts = append(parts, m.ToDomain())
	}
	return parts, rows.Err()
}

func (r *partRepository) Update(ctx context.Context, p *part.Part) error {
	query := `UPDATE parts SET name=$1, description=$2, unit=$3, price=$4, stock=$5, updated_at=$6 WHERE id=$7`
	_, err := r.db.Exec(ctx, query, p.Name(), p.Description(), p.Unit(), p.Price(), p.Stock(), p.UpdatedAt(), p.ID())
	return err
}

func (r *partRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM parts WHERE id=$1`, id)
	return err
}

func (r *partRepository) UpdateStock(ctx context.Context, id string, delta int) error {
	_, err := r.db.Exec(ctx, `UPDATE parts SET stock = stock + $1 WHERE id=$2`, delta, id)
	return err
}
