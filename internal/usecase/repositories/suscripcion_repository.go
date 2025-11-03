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

type suscripcionRepository struct {
	collection *mongo.Collection
}

func NewSuscripcionRepository(db *mongo.Database) SuscripcionRepository {
	return &suscripcionRepository{
		collection: db.Collection("suscripciones"),
	}
}

func (r *suscripcionRepository) Create(ctx context.Context, suscripcion *entity.Suscripcion) error {
	if suscripcion.ID.IsZero() {
		suscripcion.ID = primitive.NewObjectID()
	}
	if suscripcion.CreadoEn.IsZero() {
		suscripcion.CreadoEn = time.Now()
	}

	_, err := r.collection.InsertOne(ctx, suscripcion)
	return err
}

func (r *suscripcionRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*entity.Suscripcion, error) {
	var suscripcion entity.Suscripcion
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&suscripcion)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("suscripción no encontrada")
		}
		return nil, err
	}
	return &suscripcion, nil
}

func (r *suscripcionRepository) GetByUserID(ctx context.Context, userID primitive.ObjectID, limit, offset int) ([]*entity.Suscripcion, error) {
	filter := bson.M{"usuario_id": userID}

	opts := options.Find()
	opts.SetSort(bson.M{"creado_en": -1})
	opts.SetLimit(int64(limit))
	opts.SetSkip(int64(offset))

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var suscripciones []*entity.Suscripcion
	for cursor.Next(ctx) {
		var suscripcion entity.Suscripcion
		if err := cursor.Decode(&suscripcion); err != nil {
			continue
		}
		suscripciones = append(suscripciones, &suscripcion)
	}

	return suscripciones, cursor.Err()
}

func (r *suscripcionRepository) GetAll(ctx context.Context, filters map[string]interface{}, limit, offset int) ([]*entity.Suscripcion, error) {
	filter := bson.M{}
	for k, v := range filters {
		filter[k] = v
	}

	opts := options.Find()
	opts.SetSort(bson.M{"creado_en": -1})
	opts.SetLimit(int64(limit))
	opts.SetSkip(int64(offset))

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var suscripciones []*entity.Suscripcion
	for cursor.Next(ctx) {
		var suscripcion entity.Suscripcion
		if err := cursor.Decode(&suscripcion); err != nil {
			continue
		}
		suscripciones = append(suscripciones, &suscripcion)
	}

	return suscripciones, cursor.Err()
}

func (r *suscripcionRepository) Update(ctx context.Context, id primitive.ObjectID, updates map[string]interface{}) error {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": updates}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("suscripción no encontrada")
	}

	return nil
}

func (r *suscripcionRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}

	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("suscripción no encontrada")
	}

	return nil
}

func (r *suscripcionRepository) Count(ctx context.Context, filters map[string]interface{}) (int64, error) {
	filter := bson.M{}
	for k, v := range filters {
		filter[k] = v
	}

	return r.collection.CountDocuments(ctx, filter)
}

func (r *suscripcionRepository) GetActiveSuscripcionByUserID(ctx context.Context, userID primitive.ObjectID) (*entity.Suscripcion, error) {
	filter := bson.M{
		"usuario_id": userID,
		"estado":     entity.EstadoSuscripcionActiva,
	}

	var suscripcion entity.Suscripcion
	err := r.collection.FindOne(ctx, filter).Decode(&suscripcion)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil // No hay suscripción activa
		}
		return nil, err
	}
	return &suscripcion, nil
}

func (r *suscripcionRepository) CountActiveSuscripcionesByPlan(ctx context.Context, planID primitive.ObjectID) (int64, error) {
	filter := bson.M{
		"plan_id": planID,
		"estado":  entity.EstadoSuscripcionActiva,
	}

	return r.collection.CountDocuments(ctx, filter)
}

func (r *suscripcionRepository) GetSuscripcionesWithDetails(ctx context.Context, filters map[string]interface{}, limit, offset int) ([]map[string]interface{}, error) {
	pipeline := []bson.M{}

	if len(filters) > 0 {
		pipeline = append(pipeline, bson.M{"$match": filters})
	}

	pipeline = append(pipeline, bson.M{
		"$lookup": bson.M{
			"from":         "usuarios",
			"localField":   "usuario_id",
			"foreignField": "_id",
			"as":           "usuario",
		},
	})

	pipeline = append(pipeline, bson.M{
		"$lookup": bson.M{
			"from":         "planes_suscripcion",
			"localField":   "plan_id",
			"foreignField": "_id",
			"as":           "plan",
		},
	})

	pipeline = append(pipeline, bson.M{"$unwind": "$usuario"})
	pipeline = append(pipeline, bson.M{"$unwind": "$plan"})

	pipeline = append(pipeline, bson.M{"$sort": bson.M{"creado_en": -1}})

	pipeline = append(pipeline, bson.M{"$skip": offset})
	pipeline = append(pipeline, bson.M{"$limit": limit})

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []map[string]interface{}
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}
