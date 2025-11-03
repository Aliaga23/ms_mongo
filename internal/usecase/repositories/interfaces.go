package repositories

import (
	"context"
	"sw2p2go/internal/entity"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UsuarioRepository interface {
	Create(ctx context.Context, usuario *entity.Usuario) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*entity.Usuario, error)
	GetByEmail(ctx context.Context, email string) (*entity.Usuario, error)
	GetAll(ctx context.Context, filters map[string]interface{}, limit, offset int) ([]*entity.Usuario, error)
	Update(ctx context.Context, id primitive.ObjectID, updates map[string]interface{}) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	Search(ctx context.Context, query string, limit, offset int) ([]*entity.Usuario, error)
	Count(ctx context.Context, filters map[string]interface{}) (int64, error)
	EmailExists(ctx context.Context, email string, excludeID ...primitive.ObjectID) (bool, error)
}

type PlanRepository interface {
	Create(ctx context.Context, plan *entity.PlanSuscripcion) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*entity.PlanSuscripcion, error)
	GetAll(ctx context.Context, filters map[string]interface{}, limit, offset int) ([]*entity.PlanSuscripcion, error)
	Update(ctx context.Context, id primitive.ObjectID, updates map[string]interface{}) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	Count(ctx context.Context, filters map[string]interface{}) (int64, error)
	GetActivePlans(ctx context.Context, limit, offset int) ([]*entity.PlanSuscripcion, error)
}

type SuscripcionRepository interface {
	Create(ctx context.Context, suscripcion *entity.Suscripcion) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*entity.Suscripcion, error)
	GetByUserID(ctx context.Context, userID primitive.ObjectID, limit, offset int) ([]*entity.Suscripcion, error)
	GetAll(ctx context.Context, filters map[string]interface{}, limit, offset int) ([]*entity.Suscripcion, error)
	Update(ctx context.Context, id primitive.ObjectID, updates map[string]interface{}) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	Count(ctx context.Context, filters map[string]interface{}) (int64, error)
	GetActiveSuscripcionByUserID(ctx context.Context, userID primitive.ObjectID) (*entity.Suscripcion, error)
	CountActiveSuscripcionesByPlan(ctx context.Context, planID primitive.ObjectID) (int64, error)
	GetSuscripcionesWithDetails(ctx context.Context, filters map[string]interface{}, limit, offset int) ([]map[string]interface{}, error)
}
