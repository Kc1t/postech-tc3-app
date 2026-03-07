package mongodb

import (
	"context"

	mongomodel "github.com/fiap/postech-tc1/internal/adapters/outbound/mongodb/model"
	"github.com/fiap/postech-tc1/internal/domain/part"
	"github.com/fiap/postech-tc1/internal/ports"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

const partCollection = "parts"

type partRepository struct {
	collection *mongodriver.Collection
}

func NewPartRepository(db *mongodriver.Database) ports.PartRepository {
	return &partRepository{collection: db.Collection(partCollection)}
}

func (r *partRepository) Create(ctx context.Context, p *part.Part) error {
	doc := mongomodel.FromPart(p)
	_, err := r.collection.InsertOne(ctx, doc)
	if err != nil {
		return err
	}
	p.SetID(doc.ID.Hex())
	return nil
}

func (r *partRepository) FindByID(ctx context.Context, id string) (*part.Part, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var doc mongomodel.Part
	if err := r.collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&doc); err != nil {
		return nil, err
	}
	return doc.ToDomain(), nil
}

func (r *partRepository) FindAll(ctx context.Context) ([]*part.Part, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var docs []mongomodel.Part
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, err
	}

	parts := make([]*part.Part, 0, len(docs))
	for i := range docs {
		parts = append(parts, docs[i].ToDomain())
	}
	return parts, nil
}

func (r *partRepository) Update(ctx context.Context, p *part.Part) error {
	doc := mongomodel.FromPart(p)
	_, err := r.collection.ReplaceOne(ctx, bson.M{"_id": doc.ID}, doc)
	return err
}

func (r *partRepository) Delete(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": oid})
	return err
}

func (r *partRepository) UpdateStock(ctx context.Context, id string, delta int) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.collection.UpdateOne(ctx,
		bson.M{"_id": oid},
		bson.M{"$inc": bson.M{"stock": delta}},
	)
	return err
}
