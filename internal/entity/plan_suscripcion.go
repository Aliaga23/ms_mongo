package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PlanSuscripcion struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Nombre      string             `bson:"nombre" json:"nombre"`
	Descripcion string             `bson:"descripcion" json:"descripcion"`
	Precio      float64            `bson:"precio" json:"precio"`
	Activo      bool               `bson:"activo" json:"activo"`
	CreadoEn    time.Time          `bson:"creado_en" json:"creado_en"`
}

func (p PlanSuscripcion) GetCollectionName() string {
	return "planes_suscripcion"
}

func (p PlanSuscripcion) IsActive() bool {
	return p.Activo
}

func (p PlanSuscripcion) GetID() string {
	return p.ID.Hex()
}
