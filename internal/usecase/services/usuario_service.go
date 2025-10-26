package services

import (
	"context"
	"errors"
	"strings"
	"sw2p2go/config"
	"sw2p2go/internal/dto"
	"sw2p2go/internal/entity"
	"sw2p2go/internal/usecase/repositories"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type usuarioService struct {
	userRepo repositories.UsuarioRepository
	cfg      *config.Config
}

func NewUsuarioService(userRepo repositories.UsuarioRepository, cfg *config.Config) UsuarioService {
	return &usuarioService{
		userRepo: userRepo,
		cfg:      cfg,
	}
}

func (s *usuarioService) Register(ctx context.Context, req *dto.CreateUsuarioRequest) (*dto.UsuarioDTO, error) {
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	req.Nombre = strings.TrimSpace(req.Nombre)

	exists, err := s.userRepo.EmailExists(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("el email ya está registrado")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	usuario := &entity.Usuario{
		Nombre:   req.Nombre,
		Email:    req.Email,
		Telefono: req.Telefono,
		Password: string(hashedPassword),
		Estado:   true,
	}

	if err := s.userRepo.Create(ctx, usuario); err != nil {
		return nil, err
	}

	return s.entityToDTO(usuario), nil
}

func (s *usuarioService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	usuario, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("credenciales inválidas")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(usuario.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("credenciales inválidas")
	}

	if !usuario.Estado {
		return nil, errors.New("usuario inactivo")
	}

	token, err := s.generateJWT(usuario)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		Token:   token,
		Usuario: *s.entityToDTO(usuario),
	}, nil
}

func (s *usuarioService) GetProfile(ctx context.Context, userID string) (*dto.UsuarioDTO, error) {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("ID de usuario inválido")
	}

	usuario, err := s.userRepo.GetByID(ctx, objectID)
	if err != nil {
		return nil, err
	}

	return s.entityToDTO(usuario), nil
}

func (s *usuarioService) GetAllUsers(ctx context.Context, limit, offset int) ([]*dto.UsuarioDTO, int64, error) {
	usuarios, err := s.userRepo.GetAll(ctx, nil, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.userRepo.Count(ctx, nil)
	if err != nil {
		return nil, 0, err
	}

	var dtos []*dto.UsuarioDTO
	for _, usuario := range usuarios {
		dtos = append(dtos, s.entityToDTO(usuario))
	}

	return dtos, total, nil
}

func (s *usuarioService) GetUserByID(ctx context.Context, id string) (*dto.UsuarioDTO, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("ID de usuario inválido")
	}

	usuario, err := s.userRepo.GetByID(ctx, objectID)
	if err != nil {
		return nil, err
	}

	return s.entityToDTO(usuario), nil
}

func (s *usuarioService) UpdateUser(ctx context.Context, id string, req *dto.UpdateUsuarioRequest) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("ID de usuario inválido")
	}

	updates := make(map[string]interface{})

	if req.Nombre != nil {
		*req.Nombre = strings.TrimSpace(*req.Nombre)
		updates["nombre"] = *req.Nombre
	}

	if req.Telefono != nil {
		*req.Telefono = strings.TrimSpace(*req.Telefono)
		updates["telefono"] = *req.Telefono
	}

	if req.Password != nil {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		updates["password"] = string(hashedPassword)
	}

	if len(updates) == 0 {
		return errors.New("no hay campos para actualizar")
	}

	return s.userRepo.Update(ctx, objectID, updates)
}

func (s *usuarioService) DeleteUser(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("ID de usuario inválido")
	}

	return s.userRepo.Delete(ctx, objectID)
}

func (s *usuarioService) SearchUsers(ctx context.Context, query string, limit, offset int) ([]*dto.UsuarioDTO, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		usuarios, _, err := s.GetAllUsers(ctx, limit, offset)
		return usuarios, err
	}

	usuarios, err := s.userRepo.Search(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}

	var dtos []*dto.UsuarioDTO
	for _, usuario := range usuarios {
		dtos = append(dtos, s.entityToDTO(usuario))
	}

	return dtos, nil
}

func (s *usuarioService) ChangePassword(ctx context.Context, userID string, req *dto.ChangePasswordRequest) error {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.New("ID de usuario inválido")
	}

	usuario, err := s.userRepo.GetByID(ctx, objectID)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(usuario.Password), []byte(req.CurrentPassword)); err != nil {
		return errors.New("contraseña actual incorrecta")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	updates := map[string]interface{}{
		"password": string(hashedPassword),
	}

	return s.userRepo.Update(ctx, objectID, updates)
}

func (s *usuarioService) generateJWT(usuario *entity.Usuario) (string, error) {
	claims := jwt.MapClaims{
		"user_id": usuario.ID.Hex(),
		"email":   usuario.Email,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.JWTSecret))
}

func (s *usuarioService) entityToDTO(usuario *entity.Usuario) *dto.UsuarioDTO {
	return &dto.UsuarioDTO{
		ID:       usuario.ID.Hex(),
		Nombre:   usuario.Nombre,
		Email:    usuario.Email,
		Telefono: usuario.Telefono,
		Estado:   usuario.Estado,
		CreadoEn: usuario.CreadoEn,
	}
}
