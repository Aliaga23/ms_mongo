package dto

import "time"

// UsuarioDTO - DTO para respuesta de usuario
type UsuarioDTO struct {
	ID       string    `json:"id"`
	Nombre   string    `json:"nombre"`
	Email    string    `json:"email"`
	Telefono string    `json:"telefono"`
	Estado   bool      `json:"estado"`
	CreadoEn time.Time `json:"creado_en"`
}

// CreateUsuarioRequest - DTO para crear usuario
type CreateUsuarioRequest struct {
	Nombre   string `json:"nombre" binding:"required,min=2,max=100"`
	Email    string `json:"email" binding:"required,email"`
	Telefono string `json:"telefono" binding:"max=20"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginRequest - DTO para login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// UpdateUsuarioRequest - DTO para actualizar usuario
type UpdateUsuarioRequest struct {
	Nombre   *string `json:"nombre,omitempty" binding:"omitempty,min=2,max=100"`
	Telefono *string `json:"telefono,omitempty" binding:"omitempty,max=20"`
	Password *string `json:"password,omitempty" binding:"omitempty,min=6"`
}

// LoginResponse - DTO para respuesta de login
type LoginResponse struct {
	Token   string     `json:"token"`
	Usuario UsuarioDTO `json:"usuario"`
}

// ChangePasswordRequest - DTO para cambiar contraseña
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=6"`
}