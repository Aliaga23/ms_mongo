package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Usuario struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Nombre   string             `bson:"nombre" json:"nombre"`
	Email    string             `bson:"email" json:"email"`
	Telefono string             `bson:"telefono" json:"telefono"`
	Password string             `bson:"password" json:"-"`
	Estado   bool               `bson:"estado" json:"estado"`
	EsAdmin  bool               `bson:"es_admin" json:"es_admin"`
	CreadoEn time.Time          `bson:"creado_en" json:"creado_en"`
}

func (u Usuario) GetCollectionName() string {
	return "usuarios"
}

func (u Usuario) IsActive() bool {
	return u.Estado
}

func (u Usuario) GetID() string {
	return u.ID.Hex()
}
