package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Usuario representa la entidad de usuario en la base de datos
type Usuario struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Nombre   string             `bson:"nombre" json:"nombre"`
	Email    string             `bson:"email" json:"email"`
	Telefono string             `bson:"telefono" json:"telefono"`
	Password string             `bson:"password" json:"-"` // No incluir en JSON
	Estado   bool               `bson:"estado" json:"estado"`
	CreadoEn time.Time          `bson:"creado_en" json:"creado_en"`
}

// GetCollectionName retorna el nombre de la colección
func (u Usuario) GetCollectionName() string {
	return "usuarios"
}

// IsActive verifica si el usuario está activo
func (u Usuario) IsActive() bool {
	return u.Estado
}

// GetID retorna el ID del usuario como string
func (u Usuario) GetID() string {
	return u.ID.Hex()
}