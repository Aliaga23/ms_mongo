package services

import (
	"context"
	"errors"
	"strings"
	"sw2p2go/internal/dto"
	"sw2p2go/internal/entity"
	"sw2p2go/internal/usecase/repositories"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type planService struct {
	planRepo        repositories.PlanRepository
	suscripcionRepo repositories.SuscripcionRepository
}

// NewPlanService crea una nueva instancia del servicio de planes
func NewPlanService(planRepo repositories.PlanRepository, suscripcionRepo repositories.SuscripcionRepository) PlanService {
	return &planService{
		planRepo:        planRepo,
		suscripcionRepo: suscripcionRepo,
	}
}

// CreatePlan crea un nuevo plan
func (s *planService) CreatePlan(ctx context.Context, req *dto.CreatePlanRequest) (*dto.PlanSuscripcionDTO, error) {
	// Normalizar datos
	req.Nombre = strings.TrimSpace(req.Nombre)
	req.Descripcion = strings.TrimSpace(req.Descripcion)

	// Crear entidad plan
	plan := &entity.PlanSuscripcion{
		Nombre:      req.Nombre,
		Descripcion: req.Descripcion,
		Precio:      req.Precio,
		Activo:      true,
		CreadoEn:    time.Now(),
	}

	// Guardar en base de datos
	if err := s.planRepo.Create(ctx, plan); err != nil {
		return nil, err
	}

	// Convertir a DTO
	return s.entityToDTO(plan), nil
}

// GetAllPlans obtiene todos los planes con filtros y paginación
func (s *planService) GetAllPlans(ctx context.Context, showInactive bool, limit, offset int) ([]*dto.PlanSuscripcionDTO, int64, error) {
	filters := make(map[string]interface{})
	if !showInactive {
		filters["activo"] = true
	}

	planes, err := s.planRepo.GetAll(ctx, filters, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	// Contar total
	total, err := s.planRepo.Count(ctx, filters)
	if err != nil {
		return nil, 0, err
	}

	// Convertir a DTOs
	var dtos []*dto.PlanSuscripcionDTO
	for _, plan := range planes {
		dtos = append(dtos, s.entityToDTO(plan))
	}

	return dtos, total, nil
}

// GetPlanByID obtiene un plan por ID
func (s *planService) GetPlanByID(ctx context.Context, id string) (*dto.PlanSuscripcionDTO, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("ID de plan inválido")
	}

	plan, err := s.planRepo.GetByID(ctx, objectID)
	if err != nil {
		return nil, err
	}

	return s.entityToDTO(plan), nil
}

// UpdatePlan actualiza un plan
func (s *planService) UpdatePlan(ctx context.Context, id string, req *dto.UpdatePlanRequest) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("ID de plan inválido")
	}

	updates := make(map[string]interface{})

	if req.Nombre != nil {
		*req.Nombre = strings.TrimSpace(*req.Nombre)
		updates["nombre"] = *req.Nombre
	}

	if req.Descripcion != nil {
		*req.Descripcion = strings.TrimSpace(*req.Descripcion)
		updates["descripcion"] = *req.Descripcion
	}

	if req.Precio != nil {
		updates["precio"] = *req.Precio
	}

	if req.Activo != nil {
		updates["activo"] = *req.Activo
	}

	if len(updates) == 0 {
		return errors.New("no hay campos para actualizar")
	}

	return s.planRepo.Update(ctx, objectID, updates)
}

// DeletePlan elimina un plan (soft delete)
func (s *planService) DeletePlan(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("ID de plan inválido")
	}

	// Verificar si hay suscripciones activas
	activeCount, err := s.suscripcionRepo.CountActiveSuscripcionesByPlan(ctx, objectID)
	if err != nil {
		return err
	}

	if activeCount > 0 {
		return errors.New("no se puede eliminar un plan con suscripciones activas")
	}

	return s.planRepo.Delete(ctx, objectID)
}

// GetActivePlans obtiene solo los planes activos
func (s *planService) GetActivePlans(ctx context.Context, limit, offset int) ([]*dto.PlanSuscripcionDTO, error) {
	planes, err := s.planRepo.GetActivePlans(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	// Convertir a DTOs
	var dtos []*dto.PlanSuscripcionDTO
	for _, plan := range planes {
		dtos = append(dtos, s.entityToDTO(plan))
	}

	return dtos, nil
}

// entityToDTO convierte una entidad PlanSuscripcion a DTO
func (s *planService) entityToDTO(plan *entity.PlanSuscripcion) *dto.PlanSuscripcionDTO {
	return &dto.PlanSuscripcionDTO{
		ID:          plan.ID.Hex(),
		Nombre:      plan.Nombre,
		Descripcion: plan.Descripcion,
		Precio:      plan.Precio,
		Activo:      plan.Activo,
		CreadoEn:    plan.CreadoEn,
	}
}
