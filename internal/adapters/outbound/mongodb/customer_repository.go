package mongodb

import (
	"context"

	"github.com/fiap/postech-tc1/internal/adapters/outbound/mongodb/model"
	"github.com/fiap/postech-tc1/internal/domain/customer"
	"github.com/fiap/postech-tc1/internal/ports"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

const customerCollection = "customers"

type customerRepository struct {
	collection *mongodriver.Collection
}

func NewCustomerRepository(db *mongodriver.Database) ports.CustomerRepository {
	return &customerRepository{collection: db.Collection(customerCollection)}
}

func (r *customerRepository) Create(ctx context.Context, c *customer.Customer) error {
	doc := mongomodel.FromCustomer(c)
	_, err := r.collection.InsertOne(ctx, doc)
	if err != nil {
		return err
	}
	c.SetID(doc.ID.Hex())
	return nil
}

func (r *customerRepository) FindByID(ctx context.Context, id string) (*customer.Customer, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var doc mongomodel.Customer
	if err := r.collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&doc); err != nil {
		return nil, err
	}
	return doc.ToDomain(), nil
}

func (r *customerRepository) FindByDocument(ctx context.Context, document string) (*customer.Customer, error) {
	var doc mongomodel.Customer
	if err := r.collection.FindOne(ctx, bson.M{"document": document}).Decode(&doc); err != nil {
		return nil, err
	}
	return doc.ToDomain(), nil
}

func (r *customerRepository) FindAll(ctx context.Context) ([]*customer.Customer, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var docs []mongomodel.Customer
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, err
	}

	customers := make([]*customer.Customer, 0, len(docs))
	for i := range docs {
		customers = append(customers, docs[i].ToDomain())
	}
	return customers, nil
}

func (r *customerRepository) Update(ctx context.Context, c *customer.Customer) error {
	doc := mongomodel.FromCustomer(c)
	_, err := r.collection.ReplaceOne(ctx, bson.M{"_id": doc.ID}, doc)
	return err
}

func (r *customerRepository) Delete(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": oid})
	return err
}
