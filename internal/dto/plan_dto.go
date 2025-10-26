package dto

import "time"

// PlanSuscripcionDTO - DTO para respuesta de plan
type PlanSuscripcionDTO struct {
	ID          string    `json:"id"`
	Nombre      string    `json:"nombre"`
	Descripcion string    `json:"descripcion"`
	Precio      float64   `json:"precio"`
	Activo      bool      `json:"activo"`
	CreadoEn    time.Time `json:"creado_en"`
}

// CreatePlanRequest - DTO para crear plan segun
type CreatePlanRequest struct {
	Nombre      string  `json:"nombre" binding:"required,min=2,max=100"`
	Descripcion string  `json:"descripcion" binding:"required,min=10,max=500"`
	Precio      float64 `json:"precio" binding:"required,gt=0"`
}

// UpdatePlanRequest - DTO para actualizar plan te odio mongo
type UpdatePlanRequest struct {
	Nombre      *string  `json:"nombre,omitempty" binding:"omitempty,min=2,max=100"`
	Descripcion *string  `json:"descripcion,omitempty" binding:"omitempty,min=10,max=500"`
	Precio      *float64 `json:"precio,omitempty" binding:"omitempty,gt=0"`
	Activo      *bool    `json:"activo,omitempty"`
}
