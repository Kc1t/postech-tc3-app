package mongodb

import (
	"context"

	mongomodel "github.com/fiap/postech-tc1/internal/adapters/outbound/mongodb/model"
	"github.com/fiap/postech-tc1/internal/domain/serviceorder"
	"github.com/fiap/postech-tc1/internal/ports"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

const serviceOrderCollection = "service_orders"

type serviceOrderRepository struct {
	collection *mongodriver.Collection
}

func NewServiceOrderRepository(db *mongodriver.Database) ports.ServiceOrderRepository {
	return &serviceOrderRepository{collection: db.Collection(serviceOrderCollection)}
}

func (r *serviceOrderRepository) Create(ctx context.Context, so *serviceorder.ServiceOrder) error {
	doc := mongomodel.FromServiceOrder(so)
	_, err := r.collection.InsertOne(ctx, doc)
	if err != nil {
		return err
	}
	so.SetID(doc.ID.Hex())
	return nil
}

func (r *serviceOrderRepository) FindByID(ctx context.Context, id string) (*serviceorder.ServiceOrder, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var doc mongomodel.ServiceOrder
	if err := r.collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&doc); err != nil {
		return nil, err
	}
	return doc.ToDomain(), nil
}

func (r *serviceOrderRepository) FindAll(ctx context.Context) ([]*serviceorder.ServiceOrder, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var docs []mongomodel.ServiceOrder
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, err
	}

	orders := make([]*serviceorder.ServiceOrder, 0, len(docs))
	for i := range docs {
		orders = append(orders, docs[i].ToDomain())
	}
	return orders, nil
}

func (r *serviceOrderRepository) FindByCustomerID(ctx context.Context, customerID string) ([]*serviceorder.ServiceOrder, error) {
	oid, err := primitive.ObjectIDFromHex(customerID)
	if err != nil {
		return nil, err
	}
	cursor, err := r.collection.Find(ctx, bson.M{"customer_id": oid})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var docs []mongomodel.ServiceOrder
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, err
	}

	orders := make([]*serviceorder.ServiceOrder, 0, len(docs))
	for i := range docs {
		orders = append(orders, docs[i].ToDomain())
	}
	return orders, nil
}

func (r *serviceOrderRepository) UpdateStatus(ctx context.Context, id string, status serviceorder.Status) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.collection.UpdateOne(ctx,
		bson.M{"_id": oid},
		bson.M{"$set": bson.M{"status": status}},
	)
	return err
}

func (r *serviceOrderRepository) Update(ctx context.Context, so *serviceorder.ServiceOrder) error {
	doc := mongomodel.FromServiceOrder(so)
	_, err := r.collection.ReplaceOne(ctx, bson.M{"_id": doc.ID}, doc)
	return err
}

func (r *serviceOrderRepository) Delete(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": oid})
	return err
}
