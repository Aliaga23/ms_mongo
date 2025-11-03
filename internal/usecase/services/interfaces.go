package services

import (
	"context"
	"sw2p2go/internal/dto"
)

type UsuarioService interface {
	Register(ctx context.Context, req *dto.CreateUsuarioRequest) (*dto.UsuarioDTO, error)
	Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error)
	GetProfile(ctx context.Context, userID string) (*dto.UsuarioDTO, error)
	GetAllUsers(ctx context.Context, limit, offset int) ([]*dto.UsuarioDTO, int64, error)
	GetUserByID(ctx context.Context, id string) (*dto.UsuarioDTO, error)
	UpdateUser(ctx context.Context, id string, req *dto.UpdateUsuarioRequest) error
	DeleteUser(ctx context.Context, id string) error
	SearchUsers(ctx context.Context, query string, limit, offset int) ([]*dto.UsuarioDTO, error)
	ChangePassword(ctx context.Context, userID string, req *dto.ChangePasswordRequest) error
}

type PlanService interface {
	CreatePlan(ctx context.Context, req *dto.CreatePlanRequest) (*dto.PlanSuscripcionDTO, error)
	GetAllPlans(ctx context.Context, showInactive bool, limit, offset int) ([]*dto.PlanSuscripcionDTO, int64, error)
	GetPlanByID(ctx context.Context, id string) (*dto.PlanSuscripcionDTO, error)
	UpdatePlan(ctx context.Context, id string, req *dto.UpdatePlanRequest) error
	DeletePlan(ctx context.Context, id string) error
	GetActivePlans(ctx context.Context, limit, offset int) ([]*dto.PlanSuscripcionDTO, error)
}

type SuscripcionService interface {
	CreateSuscripcion(ctx context.Context, req *dto.CreateSuscripcionRequest) (*dto.SuscripcionDTO, error)
	GetAllSuscripciones(ctx context.Context, limit, offset int) ([]*dto.SuscripcionDTO, int64, error)
	GetSuscripcionByID(ctx context.Context, id string) (*dto.SuscripcionDTO, error)
	GetSuscripcionesByUser(ctx context.Context, userID string, limit, offset int) ([]*dto.SuscripcionDTO, error)
	GetMySuscripciones(ctx context.Context, userID string, limit, offset int) ([]*dto.SuscripcionDTO, error)
	UpdateSuscripcion(ctx context.Context, id string, req *dto.UpdateSuscripcionRequest) error
	CancelSuscripcion(ctx context.Context, id string) error
	GetSuscripcionesWithDetails(ctx context.Context, limit, offset int) ([]map[string]interface{}, int64, error)
}
