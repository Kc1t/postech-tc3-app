package postgresql

import (
	"context"

	pgmodel "github.com/fiap/postech-tc1/internal/adapters/outbound/postgresql/model"
	"github.com/fiap/postech-tc1/internal/domain/serviceorder"
	"github.com/fiap/postech-tc1/internal/ports"
	"github.com/jackc/pgx/v5/pgxpool"
)

type serviceOrderRepository struct {
	db *pgxpool.Pool
}

func NewServiceOrderRepository(db *pgxpool.Pool) ports.ServiceOrderRepository {
	return &serviceOrderRepository{db: db}
}

func (r *serviceOrderRepository) Create(ctx context.Context, so *serviceorder.ServiceOrder) error {
	doc := pgmodel.FromServiceOrder(so)
	query := `
		INSERT INTO service_orders (customer_id, vehicle_id, status, services, parts, total_amount, notes, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id`
	var id string
	err := r.db.QueryRow(ctx, query,
		doc.CustomerID, doc.VehicleID, doc.Status, doc.ServicesRaw, doc.PartsRaw,
		doc.TotalAmount, doc.Notes, doc.CreatedAt, doc.UpdatedAt,
	).Scan(&id)
	if err != nil {
		return err
	}
	so.SetID(id)
	return nil
}

func (r *serviceOrderRepository) FindByID(ctx context.Context, id string) (*serviceorder.ServiceOrder, error) {
	query := `SELECT id, customer_id, vehicle_id, status, services, parts, total_amount, notes, created_at, updated_at FROM service_orders WHERE id=$1`
	var m pgmodel.ServiceOrder
	err := r.db.QueryRow(ctx, query, id).Scan(
		&m.ID, &m.CustomerID, &m.VehicleID, &m.Status, &m.ServicesRaw, &m.PartsRaw,
		&m.TotalAmount, &m.Notes, &m.CreatedAt, &m.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return m.ToDomain(), nil
}

func (r *serviceOrderRepository) FindAll(ctx context.Context) ([]*serviceorder.ServiceOrder, error) {
	query := `SELECT id, customer_id, vehicle_id, status, services, parts, total_amount, notes, created_at, updated_at FROM service_orders`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*serviceorder.ServiceOrder
	for rows.Next() {
		var m pgmodel.ServiceOrder
		if err := rows.Scan(&m.ID, &m.CustomerID, &m.VehicleID, &m.Status, &m.ServicesRaw, &m.PartsRaw, &m.TotalAmount, &m.Notes, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, err
		}
		orders = append(orders, m.ToDomain())
	}
	return orders, rows.Err()
}

func (r *serviceOrderRepository) FindByCustomerID(ctx context.Context, customerID string) ([]*serviceorder.ServiceOrder, error) {
	query := `SELECT id, customer_id, vehicle_id, status, services, parts, total_amount, notes, created_at, updated_at FROM service_orders WHERE customer_id=$1`
	rows, err := r.db.Query(ctx, query, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*serviceorder.ServiceOrder
	for rows.Next() {
		var m pgmodel.ServiceOrder
		if err := rows.Scan(&m.ID, &m.CustomerID, &m.VehicleID, &m.Status, &m.ServicesRaw, &m.PartsRaw, &m.TotalAmount, &m.Notes, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, err
		}
		orders = append(orders, m.ToDomain())
	}
	return orders, rows.Err()
}

func (r *serviceOrderRepository) UpdateStatus(ctx context.Context, id string, status serviceorder.Status) error {
	_, err := r.db.Exec(ctx, `UPDATE service_orders SET status=$1 WHERE id=$2`, status, id)
	return err
}

func (r *serviceOrderRepository) Update(ctx context.Context, so *serviceorder.ServiceOrder) error {
	doc := pgmodel.FromServiceOrder(so)
	query := `UPDATE service_orders SET status=$1, services=$2, parts=$3, total_amount=$4, notes=$5, updated_at=$6 WHERE id=$7`
	_, err := r.db.Exec(ctx, query, doc.Status, doc.ServicesRaw, doc.PartsRaw, doc.TotalAmount, doc.Notes, doc.UpdatedAt, doc.ID)
	return err
}

func (r *serviceOrderRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM service_orders WHERE id=$1`, id)
	return err
}
