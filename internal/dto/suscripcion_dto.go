package dto

import "time"

// SuscripcionDTO - DTO para respuesta de suscripción
type SuscripcionDTO struct {
	ID          string              `json:"id"`
	UsuarioID   string              `json:"usuario_id"`
	PlanID      string              `json:"plan_id"`
	FechaInicio time.Time           `json:"fecha_inicio"`
	FechaFin    time.Time           `json:"fecha_fin"`
	Estado      string              `json:"estado"`
	CreadoEn    time.Time           `json:"creado_en"`
	Usuario     *UsuarioDTO         `json:"usuario,omitempty"`
	Plan        *PlanSuscripcionDTO `json:"plan,omitempty"`
}

// CreateSuscripcionRequest - DTO para crear suscripción
type CreateSuscripcionRequest struct {
	UsuarioID   string `json:"usuario_id" binding:"required"`
	PlanID      string `json:"plan_id" binding:"required"`
	FechaInicio string `json:"fecha_inicio,omitempty"`
	FechaFin    string `json:"fecha_fin,omitempty"`
}

// UpdateSuscripcionRequest - DTO para actualizar suscripción
type UpdateSuscripcionRequest struct {
	FechaFin *string `json:"fecha_fin,omitempty"`
	Estado   *string `json:"estado,omitempty" binding:"omitempty,oneof=activa vencida cancelada"`
}
