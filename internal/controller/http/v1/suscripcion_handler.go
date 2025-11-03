package v1

import (
	"net/http"
	"strconv"
	"sw2p2go/internal/dto"
	"sw2p2go/internal/usecase/services"

	"github.com/gin-gonic/gin"
)

type SuscripcionHandler struct {
	suscripcionService services.SuscripcionService
}

func NewSuscripcionHandler(suscripcionService services.SuscripcionService) *SuscripcionHandler {
	return &SuscripcionHandler{
		suscripcionService: suscripcionService,
	}
}

func (h *SuscripcionHandler) CreateSuscripcion(c *gin.Context) {
	var req dto.CreateSuscripcionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("Datos inválidos", err.Error()))
		return
	}

	suscripcion, err := h.suscripcionService.CreateSuscripcion(c.Request.Context(), &req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "el usuario ya tiene una suscripción activa" ||
			err.Error() == "usuario no encontrado" ||
			err.Error() == "plan no encontrado" ||
			err.Error() == "usuario inactivo" ||
			err.Error() == "plan inactivo" {
			statusCode = http.StatusConflict
		}
		c.JSON(statusCode, dto.NewErrorResponse("Error creando suscripción", err.Error()))
		return
	}

	c.JSON(http.StatusCreated, dto.NewSuccessResponse("Suscripción creada exitosamente", suscripcion))
}

func (h *SuscripcionHandler) GetAllSuscripciones(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	suscripciones, total, err := h.suscripcionService.GetAllSuscripciones(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse("Error obteniendo suscripciones", err.Error()))
		return
	}

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
		Message: "Suscripciones obtenidas exitosamente",
		Data:    suscripciones,
		Meta:    meta,
	}

	c.JSON(http.StatusOK, response)
}

func (h *SuscripcionHandler) GetSuscripcionByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("ID requerido", "missing_id"))
		return
	}

	suscripcion, err := h.suscripcionService.GetSuscripcionByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.NewErrorResponse("Suscripción no encontrada", err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.NewSuccessResponse("Suscripción obtenida exitosamente", suscripcion))
}

func (h *SuscripcionHandler) GetSuscripcionesByUser(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("ID de usuario requerido", "missing_user_id"))
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

	suscripciones, err := h.suscripcionService.GetSuscripcionesByUser(c.Request.Context(), userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse("Error obteniendo suscripciones del usuario", err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.NewSuccessResponse("Suscripciones del usuario obtenidas exitosamente", suscripciones))
}

func (h *SuscripcionHandler) GetMySuscripciones(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, dto.NewErrorResponse("Usuario no autenticado", "unauthorized"))
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

	suscripciones, err := h.suscripcionService.GetMySuscripciones(c.Request.Context(), userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse("Error obteniendo tus suscripciones", err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.NewSuccessResponse("Tus suscripciones obtenidas exitosamente", suscripciones))
}

func (h *SuscripcionHandler) UpdateSuscripcion(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("ID requerido", "missing_id"))
		return
	}

	var req dto.UpdateSuscripcionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("Datos inválidos", err.Error()))
		return
	}

	if err := h.suscripcionService.UpdateSuscripcion(c.Request.Context(), id, &req); err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "suscripción no encontrada" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "formato de fecha inválido (use YYYY-MM-DD)" ||
			err.Error() == "estado inválido" {
			statusCode = http.StatusBadRequest
		}
		c.JSON(statusCode, dto.NewErrorResponse("Error actualizando suscripción", err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.NewSuccessResponse("Suscripción actualizada exitosamente", nil))
}

func (h *SuscripcionHandler) CancelSuscripcion(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("ID requerido", "missing_id"))
		return
	}

	if err := h.suscripcionService.CancelSuscripcion(c.Request.Context(), id); err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "suscripción no encontrada" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, dto.NewErrorResponse("Error cancelando suscripción", err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.NewSuccessResponse("Suscripción cancelada exitosamente", nil))
}

func (h *SuscripcionHandler) GetSuscripcionesWithDetails(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	suscripciones, total, err := h.suscripcionService.GetSuscripcionesWithDetails(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse("Error obteniendo suscripciones con detalles", err.Error()))
		return
	}

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
		Message: "Suscripciones con detalles obtenidas exitosamente",
		Data:    suscripciones,
		Meta:    meta,
	}

	c.JSON(http.StatusOK, response)
}
