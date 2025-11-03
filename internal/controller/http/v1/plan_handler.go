package v1

import (
	"net/http"
	"strconv"
	"sw2p2go/internal/dto"
	"sw2p2go/internal/usecase/services"

	"github.com/gin-gonic/gin"
)

type PlanHandler struct {
	planService services.PlanService
}

func NewPlanHandler(planService services.PlanService) *PlanHandler {
	return &PlanHandler{
		planService: planService,
	}
}

func (h *PlanHandler) CreatePlan(c *gin.Context) {
	var req dto.CreatePlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("Datos inválidos", err.Error()))
		return
	}

	plan, err := h.planService.CreatePlan(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse("Error creando plan", err.Error()))
		return
	}

	c.JSON(http.StatusCreated, dto.NewSuccessResponse("Plan creado exitosamente", plan))
}

func (h *PlanHandler) GetAllPlanes(c *gin.Context) {
	showInactive, _ := strconv.ParseBool(c.DefaultQuery("show_inactive", "false"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	planes, total, err := h.planService.GetAllPlans(c.Request.Context(), showInactive, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse("Error obteniendo planes", err.Error()))
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
		Message: "Planes obtenidos exitosamente",
		Data:    planes,
		Meta:    meta,
	}

	c.JSON(http.StatusOK, response)
}

func (h *PlanHandler) GetPlanByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("ID requerido", "missing_id"))
		return
	}

	plan, err := h.planService.GetPlanByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.NewErrorResponse("Plan no encontrado", err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.NewSuccessResponse("Plan obtenido exitosamente", plan))
}

func (h *PlanHandler) UpdatePlan(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("ID requerido", "missing_id"))
		return
	}

	var req dto.UpdatePlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("Datos inválidos", err.Error()))
		return
	}

	if err := h.planService.UpdatePlan(c.Request.Context(), id, &req); err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "plan no encontrado" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, dto.NewErrorResponse("Error actualizando plan", err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.NewSuccessResponse("Plan actualizado exitosamente", nil))
}

func (h *PlanHandler) DeletePlan(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("ID requerido", "missing_id"))
		return
	}

	if err := h.planService.DeletePlan(c.Request.Context(), id); err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "plan no encontrado" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "no se puede eliminar un plan con suscripciones activas" {
			statusCode = http.StatusConflict
		}
		c.JSON(statusCode, dto.NewErrorResponse("Error eliminando plan", err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.NewSuccessResponse("Plan eliminado exitosamente", nil))
}

func (h *PlanHandler) GetActivePlanes(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	planes, err := h.planService.GetActivePlans(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse("Error obteniendo planes activos", err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.NewSuccessResponse("Planes activos obtenidos exitosamente", planes))
}
