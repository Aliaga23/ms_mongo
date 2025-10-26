package v1

import (
	"net/http"
	"strconv"
	"sw2p2go/internal/dto"
	"sw2p2go/internal/usecase/services"

	"github.com/gin-gonic/gin"
)

type UsuarioHandler struct {
	usuarioService services.UsuarioService
}

func NewUsuarioHandler(usuarioService services.UsuarioService) *UsuarioHandler {
	return &UsuarioHandler{
		usuarioService: usuarioService,
	}
}

func (h *UsuarioHandler) Register(c *gin.Context) {
	var req dto.CreateUsuarioRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("Datos inválidos", err.Error()))
		return
	}

	usuario, err := h.usuarioService.Register(c.Request.Context(), &req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "el email ya está registrado" {
			statusCode = http.StatusConflict
		}
		c.JSON(statusCode, dto.NewErrorResponse("Error en el registro", err.Error()))
		return
	}

	c.JSON(http.StatusCreated, dto.NewSuccessResponse("Usuario registrado exitosamente", usuario))
}

func (h *UsuarioHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("Datos inválidos", err.Error()))
		return
	}

	response, err := h.usuarioService.Login(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.NewErrorResponse("Error en el login", err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.NewSuccessResponse("Login exitoso", response))
}

func (h *UsuarioHandler) GetProfile(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, dto.NewErrorResponse("Usuario no autenticado", "unauthorized"))
		return
	}

	usuario, err := h.usuarioService.GetProfile(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.NewErrorResponse("Usuario no encontrado", err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.NewSuccessResponse("Perfil obtenido exitosamente", usuario))
}

func (h *UsuarioHandler) GetAllUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	usuarios, total, err := h.usuarioService.GetAllUsers(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse("Error obteniendo usuarios", err.Error()))
		return
	}

	// Calcular metadatos de paginación
	totalPages := int((total + int64(limit) - 1) / int64(limit))
	meta := dto.MetaData{
		Page:        page,
		Limit:       limit,
		Total:       total,
		TotalPages:  totalPages,
		HasNext:     page < totalPages,
		HasPrevious: page > 1,
	}

	response := &dto.PaginatedResponse{
		Success: true,
		Message: "Usuarios obtenidos exitosamente",
		Data:    usuarios,
		Meta:    meta,
	}

	c.JSON(http.StatusOK, response)
}

func (h *UsuarioHandler) GetUserByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("ID requerido", "missing_id"))
		return
	}

	usuario, err := h.usuarioService.GetUserByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.NewErrorResponse("Usuario no encontrado", err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.NewSuccessResponse("Usuario obtenido exitosamente", usuario))
}

func (h *UsuarioHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("ID requerido", "missing_id"))
		return
	}

	var req dto.UpdateUsuarioRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("Datos inválidos", err.Error()))
		return
	}

	if err := h.usuarioService.UpdateUser(c.Request.Context(), id, &req); err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "usuario no encontrado" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, dto.NewErrorResponse("Error actualizando usuario", err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.NewSuccessResponse("Usuario actualizado exitosamente", nil))
}

func (h *UsuarioHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("ID requerido", "missing_id"))
		return
	}

	if err := h.usuarioService.DeleteUser(c.Request.Context(), id); err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "usuario no encontrado" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, dto.NewErrorResponse("Error eliminando usuario", err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.NewSuccessResponse("Usuario eliminado exitosamente", nil))
}

func (h *UsuarioHandler) SearchUsers(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("Parámetro de búsqueda requerido", "missing_query"))
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	usuarios, err := h.usuarioService.SearchUsers(c.Request.Context(), query, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse("Error en la búsqueda", err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.NewSuccessResponse("Búsqueda completada exitosamente", usuarios))
}

func (h *UsuarioHandler) ChangePassword(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, dto.NewErrorResponse("Usuario no autenticado", "unauthorized"))
		return
	}

	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("Datos inválidos", err.Error()))
		return
	}

	if err := h.usuarioService.ChangePassword(c.Request.Context(), userID, &req); err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "contraseña actual incorrecta" {
			statusCode = http.StatusBadRequest
		}
		c.JSON(statusCode, dto.NewErrorResponse("Error cambiando contraseña", err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.NewSuccessResponse("Contraseña cambiada exitosamente", nil))
}
