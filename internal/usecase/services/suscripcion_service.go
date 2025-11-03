package services

import (
	"context"
	"errors"
	"sw2p2go/internal/dto"
	"sw2p2go/internal/entity"
	"sw2p2go/internal/usecase/repositories"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type suscripcionService struct {
	suscripcionRepo repositories.SuscripcionRepository
	userRepo        repositories.UsuarioRepository
	planRepo        repositories.PlanRepository
}

func NewSuscripcionService(
	suscripcionRepo repositories.SuscripcionRepository,
	userRepo repositories.UsuarioRepository,
	planRepo repositories.PlanRepository,
) SuscripcionService {
	return &suscripcionService{
		suscripcionRepo: suscripcionRepo,
		userRepo:        userRepo,
		planRepo:        planRepo,
	}
}

func (s *suscripcionService) CreateSuscripcion(ctx context.Context, req *dto.CreateSuscripcionRequest) (*dto.SuscripcionDTO, error) {
	userID, err := primitive.ObjectIDFromHex(req.UsuarioID)
	if err != nil {
		return nil, errors.New("ID de usuario inválido")
	}

	planID, err := primitive.ObjectIDFromHex(req.PlanID)
	if err != nil {
		return nil, errors.New("ID de plan inválido")
	}

	usuario, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.New("usuario no encontrado")
	}
	if !usuario.Estado {
		return nil, errors.New("usuario inactivo")
	}

	plan, err := s.planRepo.GetByID(ctx, planID)
	if err != nil {
		return nil, errors.New("plan no encontrado")
	}
	if !plan.Activo {
		return nil, errors.New("plan inactivo")
	}

	activeSuscripcion, err := s.suscripcionRepo.GetActiveSuscripcionByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if activeSuscripcion != nil {
		return nil, errors.New("el usuario ya tiene una suscripción activa")
	}

	fechaInicio := time.Now()
	if req.FechaInicio != "" {
		if parsedTime, err := time.Parse("2006-01-02", req.FechaInicio); err == nil {
			fechaInicio = parsedTime
		}
	}

	fechaFin := fechaInicio.AddDate(0, 1, 0)
	if req.FechaFin != "" {
		if parsedTime, err := time.Parse("2006-01-02", req.FechaFin); err == nil {
			fechaFin = parsedTime
		}
	}

	suscripcion := &entity.Suscripcion{
		UsuarioID:   userID,
		PlanID:      planID,
		FechaInicio: fechaInicio,
		FechaFin:    fechaFin,
		Estado:      entity.EstadoSuscripcionActiva,
		CreadoEn:    time.Now(),
	}

	if err := s.suscripcionRepo.Create(ctx, suscripcion); err != nil {
		return nil, err
	}

	dtoResult := s.entityToDTO(suscripcion)
	dtoResult.Usuario = &dto.UsuarioDTO{
		ID:       usuario.ID.Hex(),
		Nombre:   usuario.Nombre,
		Email:    usuario.Email,
		Telefono: usuario.Telefono,
		Estado:   usuario.Estado,
		CreadoEn: usuario.CreadoEn,
	}
	dtoResult.Plan = &dto.PlanSuscripcionDTO{
		ID:          plan.ID.Hex(),
		Nombre:      plan.Nombre,
		Descripcion: plan.Descripcion,
		Precio:      plan.Precio,
		Activo:      plan.Activo,
		CreadoEn:    plan.CreadoEn,
	}

	return dtoResult, nil
}

func (s *suscripcionService) GetAllSuscripciones(ctx context.Context, limit, offset int) ([]*dto.SuscripcionDTO, int64, error) {
	suscripciones, err := s.suscripcionRepo.GetAll(ctx, nil, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.suscripcionRepo.Count(ctx, nil)
	if err != nil {
		return nil, 0, err
	}

	var dtos []*dto.SuscripcionDTO
	for _, suscripcion := range suscripciones {
		dtos = append(dtos, s.entityToDTO(suscripcion))
	}

	return dtos, total, nil
}

func (s *suscripcionService) GetSuscripcionByID(ctx context.Context, id string) (*dto.SuscripcionDTO, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("ID de suscripción inválido")
	}

	suscripcion, err := s.suscripcionRepo.GetByID(ctx, objectID)
	if err != nil {
		return nil, err
	}

	return s.entityToDTO(suscripcion), nil
}

func (s *suscripcionService) GetSuscripcionesByUser(ctx context.Context, userID string, limit, offset int) ([]*dto.SuscripcionDTO, error) {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("ID de usuario inválido")
	}

	suscripciones, err := s.suscripcionRepo.GetByUserID(ctx, objectID, limit, offset)
	if err != nil {
		return nil, err
	}

	var dtos []*dto.SuscripcionDTO
	for _, suscripcion := range suscripciones {
		dtos = append(dtos, s.entityToDTO(suscripcion))
	}

	return dtos, nil
}

func (s *suscripcionService) GetMySuscripciones(ctx context.Context, userID string, limit, offset int) ([]*dto.SuscripcionDTO, error) {
	return s.GetSuscripcionesByUser(ctx, userID, limit, offset)
}

func (s *suscripcionService) UpdateSuscripcion(ctx context.Context, id string, req *dto.UpdateSuscripcionRequest) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("ID de suscripción inválido")
	}

	updates := make(map[string]interface{})

	if req.FechaFin != nil {
		if parsedTime, err := time.Parse("2006-01-02", *req.FechaFin); err == nil {
			updates["fecha_fin"] = parsedTime
		} else {
			return errors.New("formato de fecha inválido (use YYYY-MM-DD)")
		}
	}

	if req.Estado != nil {
		if *req.Estado != entity.EstadoSuscripcionActiva &&
			*req.Estado != entity.EstadoSuscripcionVencida &&
			*req.Estado != entity.EstadoSuscripcionCancelada {
			return errors.New("estado inválido")
		}
		updates["estado"] = *req.Estado
	}

	if len(updates) == 0 {
		return errors.New("no hay campos para actualizar")
	}

	return s.suscripcionRepo.Update(ctx, objectID, updates)
}

func (s *suscripcionService) CancelSuscripcion(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("ID de suscripción inválido")
	}

	updates := map[string]interface{}{
		"estado":    entity.EstadoSuscripcionCancelada,
		"fecha_fin": time.Now(),
	}

	return s.suscripcionRepo.Update(ctx, objectID, updates)
}

func (s *suscripcionService) GetSuscripcionesWithDetails(ctx context.Context, limit, offset int) ([]map[string]interface{}, int64, error) {
	results, err := s.suscripcionRepo.GetSuscripcionesWithDetails(ctx, nil, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.suscripcionRepo.Count(ctx, nil)
	if err != nil {
		return nil, 0, err
	}

	return results, total, nil
}

func (s *suscripcionService) entityToDTO(suscripcion *entity.Suscripcion) *dto.SuscripcionDTO {
	return &dto.SuscripcionDTO{
		ID:          suscripcion.ID.Hex(),
		UsuarioID:   suscripcion.UsuarioID.Hex(),
		PlanID:      suscripcion.PlanID.Hex(),
		FechaInicio: suscripcion.FechaInicio,
		FechaFin:    suscripcion.FechaFin,
		Estado:      suscripcion.Estado,
		CreadoEn:    suscripcion.CreadoEn,
	}
}
