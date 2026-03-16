package postgresql

import (
	"context"

	pgmodel "github.com/fiap/postech-tc1/internal/adapters/outbound/postgresql/model"
	"github.com/fiap/postech-tc1/internal/domain/vehicle"
	"github.com/fiap/postech-tc1/internal/ports"
	"github.com/jackc/pgx/v5/pgxpool"
)

type vehicleRepository struct {
	db *pgxpool.Pool
}

func NewVehicleRepository(db *pgxpool.Pool) ports.VehicleRepository {
	return &vehicleRepository{db: db}
}

func (r *vehicleRepository) Create(ctx context.Context, v *vehicle.Vehicle) error {
	query := `
		INSERT INTO vehicles (customer_id, plate, brand, model, year, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id`
	var id string
	err := r.db.QueryRow(ctx, query, v.CustomerID(), v.Plate(), v.Brand(), v.Model(), v.Year(), v.CreatedAt(), v.UpdatedAt()).Scan(&id)
	if err != nil {
		return err
	}
	v.SetID(id)
	return nil
}

func (r *vehicleRepository) FindByID(ctx context.Context, id string) (*vehicle.Vehicle, error) {
	query := `SELECT id, customer_id, plate, brand, model, year, created_at, updated_at FROM vehicles WHERE id=$1`
	var m pgmodel.Vehicle
	err := r.db.QueryRow(ctx, query, id).Scan(&m.ID, &m.CustomerID, &m.Plate, &m.Brand, &m.Model, &m.Year, &m.CreatedAt, &m.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return m.ToDomain(), nil
}

func (r *vehicleRepository) FindByCustomerID(ctx context.Context, customerID string) ([]*vehicle.Vehicle, error) {
	query := `SELECT id, customer_id, plate, brand, model, year, created_at, updated_at FROM vehicles WHERE customer_id=$1`
	rows, err := r.db.Query(ctx, query, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var vehicles []*vehicle.Vehicle
	for rows.Next() {
		var m pgmodel.Vehicle
		if err := rows.Scan(&m.ID, &m.CustomerID, &m.Plate, &m.Brand, &m.Model, &m.Year, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, err
		}
		vehicles = append(vehicles, m.ToDomain())
	}
	return vehicles, rows.Err()
}

func (r *vehicleRepository) FindAll(ctx context.Context) ([]*vehicle.Vehicle, error) {
	query := `SELECT id, customer_id, plate, brand, model, year, created_at, updated_at FROM vehicles`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var vehicles []*vehicle.Vehicle
	for rows.Next() {
		var m pgmodel.Vehicle
		if err := rows.Scan(&m.ID, &m.CustomerID, &m.Plate, &m.Brand, &m.Model, &m.Year, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, err
		}
		vehicles = append(vehicles, m.ToDomain())
	}
	return vehicles, rows.Err()
}

func (r *vehicleRepository) Update(ctx context.Context, v *vehicle.Vehicle) error {
	query := `UPDATE vehicles SET plate=$1, brand=$2, model=$3, year=$4, updated_at=$5 WHERE id=$6`
	_, err := r.db.Exec(ctx, query, v.Plate(), v.Brand(), v.Model(), v.Year(), v.UpdatedAt(), v.ID())
	return err
}

func (r *vehicleRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM vehicles WHERE id=$1`, id)
	return err
}
