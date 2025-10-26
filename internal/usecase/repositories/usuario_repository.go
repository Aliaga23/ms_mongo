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

type usuarioRepository struct {
	collection *mongo.Collection
}

// NewUsuarioRepository crea una nueva instancia del repositorio de usuarios
func NewUsuarioRepository(db *mongo.Database) UsuarioRepository {
	return &usuarioRepository{
		collection: db.Collection("usuarios"),
	}
}

// Create crea un nuevo usuario
func (r *usuarioRepository) Create(ctx context.Context, usuario *entity.Usuario) error {
	if usuario.ID.IsZero() {
		usuario.ID = primitive.NewObjectID()
	}
	if usuario.CreadoEn.IsZero() {
		usuario.CreadoEn = time.Now()
	}

	_, err := r.collection.InsertOne(ctx, usuario)
	return err
}

// GetByID obtiene un usuario por ID
func (r *usuarioRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*entity.Usuario, error) {
	var usuario entity.Usuario
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&usuario)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("usuario no encontrado")
		}
		return nil, err
	}
	return &usuario, nil
}

// GetByEmail obtiene un usuario por email
func (r *usuarioRepository) GetByEmail(ctx context.Context, email string) (*entity.Usuario, error) {
	var usuario entity.Usuario
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&usuario)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("usuario no encontrado")
		}
		return nil, err
	}
	return &usuario, nil
}

// GetAll obtiene todos los usuarios con filtros, límite y offset
func (r *usuarioRepository) GetAll(ctx context.Context, filters map[string]interface{}, limit, offset int) ([]*entity.Usuario, error) {
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

	var usuarios []*entity.Usuario
	for cursor.Next(ctx) {
		var usuario entity.Usuario
		if err := cursor.Decode(&usuario); err != nil {
			continue
		}
		usuarios = append(usuarios, &usuario)
	}

	return usuarios, cursor.Err()
}

// Update actualiza un usuario
func (r *usuarioRepository) Update(ctx context.Context, id primitive.ObjectID, updates map[string]interface{}) error {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": updates}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("usuario no encontrado")
	}

	return nil
}

// Delete elimina un usuario (soft delete)
func (r *usuarioRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"estado": false}}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("usuario no encontrado")
	}

	return nil
}

// Search busca usuarios por nombre o email
func (r *usuarioRepository) Search(ctx context.Context, query string, limit, offset int) ([]*entity.Usuario, error) {
	filter := bson.M{
		"$or": []bson.M{
			{"nombre": bson.M{"$regex": query, "$options": "i"}},
			{"email": bson.M{"$regex": query, "$options": "i"}},
		},
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

	var usuarios []*entity.Usuario
	for cursor.Next(ctx) {
		var usuario entity.Usuario
		if err := cursor.Decode(&usuario); err != nil {
			continue
		}
		usuarios = append(usuarios, &usuario)
	}

	return usuarios, cursor.Err()
}

// Count cuenta usuarios con filtros
func (r *usuarioRepository) Count(ctx context.Context, filters map[string]interface{}) (int64, error) {
	filter := bson.M{}
	for k, v := range filters {
		filter[k] = v
	}

	return r.collection.CountDocuments(ctx, filter)
}

// EmailExists verifica si un email ya existe
func (r *usuarioRepository) EmailExists(ctx context.Context, email string, excludeID ...primitive.ObjectID) (bool, error) {
	filter := bson.M{"email": email}

	if len(excludeID) > 0 && !excludeID[0].IsZero() {
		filter["_id"] = bson.M{"$ne": excludeID[0]}
	}

	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
