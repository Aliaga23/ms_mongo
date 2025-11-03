package repositories

import (
	"context"
	"errors"
	"sw2p2go/internal/entity"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type planRepository struct {
	collection *mongo.Collection
}

func NewPlanRepository(db *mongo.Database) PlanRepository {
	return &planRepository{
		collection: db.Collection("planes_suscripcion"),
	}
}

func (r *planRepository) Create(ctx context.Context, plan *entity.PlanSuscripcion) error {
	if plan.ID.IsZero() {
		plan.ID = primitive.NewObjectID()
	}
	if plan.CreadoEn.IsZero() {
		plan.CreadoEn = time.Now()
	}

	_, err := r.collection.InsertOne(ctx, plan)
	return err
}

func (r *planRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*entity.PlanSuscripcion, error) {
	var plan entity.PlanSuscripcion
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&plan)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("plan no encontrado")
		}
		return nil, err
	}
	return &plan, nil
}

func (r *planRepository) GetAll(ctx context.Context, filters map[string]interface{}, limit, offset int) ([]*entity.PlanSuscripcion, error) {
	filter := bson.M{}
	for k, v := range filters {
		filter[k] = v
	}

	opts := options.Find()
	opts.SetSort(bson.M{"precio": 1}) // Ordenar por precio ascendente
	opts.SetLimit(int64(limit))
	opts.SetSkip(int64(offset))

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var planes []*entity.PlanSuscripcion
	for cursor.Next(ctx) {
		var plan entity.PlanSuscripcion
		if err := cursor.Decode(&plan); err != nil {
			continue
		}
		planes = append(planes, &plan)
	}

	return planes, cursor.Err()
}

func (r *planRepository) Update(ctx context.Context, id primitive.ObjectID, updates map[string]interface{}) error {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": updates}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("plan no encontrado")
	}

	return nil
}

func (r *planRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"activo": false}}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("plan no encontrado")
	}

	return nil
}

func (r *planRepository) Count(ctx context.Context, filters map[string]interface{}) (int64, error) {
	filter := bson.M{}
	for k, v := range filters {
		filter[k] = v
	}

	return r.collection.CountDocuments(ctx, filter)
}

func (r *planRepository) GetActivePlans(ctx context.Context, limit, offset int) ([]*entity.PlanSuscripcion, error) {
	filter := bson.M{"activo": true}

	opts := options.Find()
	opts.SetSort(bson.M{"precio": 1})
	opts.SetLimit(int64(limit))
	opts.SetSkip(int64(offset))

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var planes []*entity.PlanSuscripcion
	for cursor.Next(ctx) {
		var plan entity.PlanSuscripcion
		if err := cursor.Decode(&plan); err != nil {
			continue
		}
		planes = append(planes, &plan)
	}

	return planes, cursor.Err()
}
