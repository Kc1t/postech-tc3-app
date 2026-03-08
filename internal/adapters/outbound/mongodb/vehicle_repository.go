package mongodb

import (
	"context"

	mongomodel "github.com/fiap/postech-tc1/internal/adapters/outbound/mongodb/model"
	"github.com/fiap/postech-tc1/internal/domain/vehicle"
	"github.com/fiap/postech-tc1/internal/ports"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

const vehicleCollection = "vehicles"

type vehicleRepository struct {
	collection *mongodriver.Collection
}

func NewVehicleRepository(db *mongodriver.Database) ports.VehicleRepository {
	return &vehicleRepository{collection: db.Collection(vehicleCollection)}
}

func (r *vehicleRepository) Create(ctx context.Context, v *vehicle.Vehicle) error {
	doc := mongomodel.FromVehicle(v)
	_, err := r.collection.InsertOne(ctx, doc)
	if err != nil {
		return err
	}
	v.SetID(doc.ID.Hex())
	return nil
}

func (r *vehicleRepository) FindByID(ctx context.Context, id string) (*vehicle.Vehicle, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var doc mongomodel.Vehicle
	if err := r.collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&doc); err != nil {
		return nil, err
	}
	return doc.ToDomain(), nil
}

func (r *vehicleRepository) FindByCustomerID(ctx context.Context, customerID string) ([]*vehicle.Vehicle, error) {
	oid, err := primitive.ObjectIDFromHex(customerID)
	if err != nil {
		return nil, err
	}
	cursor, err := r.collection.Find(ctx, bson.M{"customer_id": oid})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var docs []mongomodel.Vehicle
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, err
	}

	vehicles := make([]*vehicle.Vehicle, 0, len(docs))
	for i := range docs {
		vehicles = append(vehicles, docs[i].ToDomain())
	}
	return vehicles, nil
}

func (r *vehicleRepository) FindAll(ctx context.Context) ([]*vehicle.Vehicle, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var docs []mongomodel.Vehicle
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, err
	}

	vehicles := make([]*vehicle.Vehicle, 0, len(docs))
	for i := range docs {
		vehicles = append(vehicles, docs[i].ToDomain())
	}
	return vehicles, nil
}

func (r *vehicleRepository) Update(ctx context.Context, v *vehicle.Vehicle) error {
	doc := mongomodel.FromVehicle(v)
	_, err := r.collection.ReplaceOne(ctx, bson.M{"_id": doc.ID}, doc)
	return err
}

func (r *vehicleRepository) Delete(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": oid})
	return err
}
