package mealhttp

import (
	"net/http"
	"strconv"

	mealusecases "macabi-back/internal/meal/application/usecases"
	sharederrors "macabi-back/internal/shared/errors"
	"macabi-back/internal/shared/pagination"

	"github.com/gin-gonic/gin"
)

type MealTemplateHandler struct {
	createMealTemplateUC  *mealusecases.CreateMealTemplate
	listMealTemplatesUC   *mealusecases.ListMealTemplates
	updateMealTemplateUC  *mealusecases.UpdateMealTemplate
	deleteMealTemplateUC  *mealusecases.DeleteMealTemplate
	addGarnishOptionUC    *mealusecases.AddGarnishOption
	removeGarnishOptionUC *mealusecases.RemoveGarnishOption
}

func NewMealTemplateHandler(
	createUC *mealusecases.CreateMealTemplate,
	listUC *mealusecases.ListMealTemplates,
	updateUC *mealusecases.UpdateMealTemplate,
	deleteUC *mealusecases.DeleteMealTemplate,
	addGarnishUC *mealusecases.AddGarnishOption,
	removeGarnishUC *mealusecases.RemoveGarnishOption,
) *MealTemplateHandler {
	return &MealTemplateHandler{
		createMealTemplateUC:  createUC,
		listMealTemplatesUC:   listUC,
		updateMealTemplateUC:  updateUC,
		deleteMealTemplateUC:  deleteUC,
		addGarnishOptionUC:    addGarnishUC,
		removeGarnishOptionUC: removeGarnishUC,
	}
}

func (h *MealTemplateHandler) Create(c *gin.Context) {
	var req CreateMealTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse(err.Error()))
		return
	}

	tmpl, err := h.createMealTemplateUC.Execute(c.Request.Context(), mealusecases.CreateMealTemplateInput{
		Title:       req.Title,
		ImageURL:    req.ImageURL,
		Description: req.Description,
		Category:    req.Category,
		Type:        req.Type,
	})
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, ToMealTemplateResponse(tmpl))
}

func (h *MealTemplateHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "50"))
	params := pagination.NewParams(page, pageSize)

	result, err := h.listMealTemplatesUC.Execute(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederrors.NewErrorResponse(err.Error()))
		return
	}

	response := PaginatedMealTemplateResponse{
		Total:      result.Total,
		Page:       result.Page,
		PageSize:   result.PageSize,
		TotalPages: result.TotalPages,
	}
	response.Data = make([]MealTemplateResponse, len(result.Data))
	for i := range result.Data {
		response.Data[i] = ToMealTemplateResponse(&result.Data[i])
	}

	c.JSON(http.StatusOK, response)
}

func (h *MealTemplateHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req UpdateMealTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse(err.Error()))
		return
	}

	tmpl, err := h.updateMealTemplateUC.Execute(c.Request.Context(), mealusecases.UpdateMealTemplateInput{
		ID:          id,
		Title:       req.Title,
		ImageURL:    req.ImageURL,
		Description: req.Description,
		Category:    req.Category,
		Type:        req.Type,
	})
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, ToMealTemplateResponse(tmpl))
}

func (h *MealTemplateHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.deleteMealTemplateUC.Execute(c.Request.Context(), id); err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "template eliminado correctamente"})
}

func (h *MealTemplateHandler) AddGarnishOption(c *gin.Context) {
	templateID := c.Param("id")
	var req CreateGarnishOptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse(err.Error()))
		return
	}
	option, err := h.addGarnishOptionUC.Execute(c.Request.Context(), mealusecases.AddGarnishOptionInput{
		TemplateID: templateID,
		Name:       req.Name,
	})
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusCreated, GarnishOptionResponse{ID: option.ID, Name: option.Name})
}

func (h *MealTemplateHandler) RemoveGarnishOption(c *gin.Context) {
	templateID := c.Param("id")
	optionID := c.Param("garnishId")
	if err := h.removeGarnishOptionUC.Execute(c.Request.Context(), mealusecases.RemoveGarnishOptionInput{
		TemplateID: templateID,
		OptionID:   optionID,
	}); err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "guarnición eliminada correctamente"})
}
