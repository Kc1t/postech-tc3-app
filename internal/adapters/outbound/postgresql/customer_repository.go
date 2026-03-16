package postgresql

import (
	"context"

	pgmodel "github.com/fiap/postech-tc1/internal/adapters/outbound/postgresql/model"
	"github.com/fiap/postech-tc1/internal/domain/customer"
	"github.com/fiap/postech-tc1/internal/ports"
	"github.com/jackc/pgx/v5/pgxpool"
)

type customerRepository struct {
	db *pgxpool.Pool
}

func NewCustomerRepository(db *pgxpool.Pool) ports.CustomerRepository {
	return &customerRepository{db: db}
}

func (r *customerRepository) Create(ctx context.Context, c *customer.Customer) error {
	query := `
		INSERT INTO customers (name, document, email, phone, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`
	var id string
	err := r.db.QueryRow(ctx, query, c.Name(), c.Document(), c.Email(), c.Phone(), c.CreatedAt(), c.UpdatedAt()).Scan(&id)
	if err != nil {
		return err
	}
	c.SetID(id)
	return nil
}

func (r *customerRepository) FindByID(ctx context.Context, id string) (*customer.Customer, error) {
	query := `SELECT id, name, document, email, phone, created_at, updated_at FROM customers WHERE id = $1`
	var m pgmodel.Customer
	err := r.db.QueryRow(ctx, query, id).Scan(&m.ID, &m.Name, &m.Document, &m.Email, &m.Phone, &m.CreatedAt, &m.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return m.ToDomain(), nil
}

func (r *customerRepository) FindByDocument(ctx context.Context, document string) (*customer.Customer, error) {
	query := `SELECT id, name, document, email, phone, created_at, updated_at FROM customers WHERE document = $1`
	var m pgmodel.Customer
	err := r.db.QueryRow(ctx, query, document).Scan(&m.ID, &m.Name, &m.Document, &m.Email, &m.Phone, &m.CreatedAt, &m.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return m.ToDomain(), nil
}

func (r *customerRepository) FindAll(ctx context.Context) ([]*customer.Customer, error) {
	query := `SELECT id, name, document, email, phone, created_at, updated_at FROM customers`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var customers []*customer.Customer
	for rows.Next() {
		var m pgmodel.Customer
		if err := rows.Scan(&m.ID, &m.Name, &m.Document, &m.Email, &m.Phone, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, err
		}
		customers = append(customers, m.ToDomain())
	}
	return customers, rows.Err()
}

func (r *customerRepository) Update(ctx context.Context, c *customer.Customer) error {
	query := `UPDATE customers SET name=$1, email=$2, phone=$3, updated_at=$4 WHERE id=$5`
	_, err := r.db.Exec(ctx, query, c.Name(), c.Email(), c.Phone(), c.UpdatedAt(), c.ID())
	return err
}

func (r *customerRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM customers WHERE id=$1`, id)
	return err
}
