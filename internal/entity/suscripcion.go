package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Suscripcion struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UsuarioID   primitive.ObjectID `bson:"usuario_id" json:"usuario_id"`
	PlanID      primitive.ObjectID `bson:"plan_id" json:"plan_id"`
	FechaInicio time.Time          `bson:"fecha_inicio" json:"fecha_inicio"`
	FechaFin    time.Time          `bson:"fecha_fin" json:"fecha_fin"`
	Estado      string             `bson:"estado" json:"estado"` // activa, vencida, cancelada
	CreadoEn    time.Time          `bson:"creado_en" json:"creado_en"`
}

const (
	EstadoSuscripcionActiva    = "activa"
	EstadoSuscripcionVencida   = "vencida"
	EstadoSuscripcionCancelada = "cancelada"
)

func (s Suscripcion) GetCollectionName() string {
	return "suscripciones"
}

func (s Suscripcion) IsActive() bool {
	return s.Estado == EstadoSuscripcionActiva && time.Now().Before(s.FechaFin)
}

func (s Suscripcion) IsExpired() bool {
	return time.Now().After(s.FechaFin)
}

func (s Suscripcion) GetID() string {
	return s.ID.Hex()
}

func (s Suscripcion) GetUsuarioID() string {
	return s.UsuarioID.Hex()
}

func (s Suscripcion) GetPlanID() string {
	return s.PlanID.Hex()
}
