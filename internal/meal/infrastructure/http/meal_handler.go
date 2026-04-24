package mealhttp

import (
	"net/http"
	"strconv"
	"time"

	mealusecases "macabi-back/internal/meal/application/usecases"
	sharederrors "macabi-back/internal/shared/errors"
	"macabi-back/internal/shared/pagination"

	"github.com/gin-gonic/gin"
)

type MealHandler struct {
	createMealUC         *mealusecases.CreateMeal
	listAvailableMealsUC *mealusecases.ListAvailableMeals
}

func NewMealHandler(createMealUC *mealusecases.CreateMeal, listAvailableMealsUC *mealusecases.ListAvailableMeals) *MealHandler {
	return &MealHandler{createMealUC: createMealUC, listAvailableMealsUC: listAvailableMealsUC}
}

func (h *MealHandler) Create(c *gin.Context) {
	var req CreateMealRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse(err.Error()))
		return
	}

	meal, err := h.createMealUC.Execute(c.Request.Context(), mealusecases.CreateMealInput{
		Title:          req.Title,
		ImageURL:       req.ImageURL,
		Description:    req.Description,
		Category:       req.Category,
		Type:           req.Type,
		AvailableCount: req.AvailableCount,
		Date:           req.Date,
	})
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, ToMealResponse(meal))
}

func (h *MealHandler) ListByDate(c *gin.Context) {
	dateStr := c.Query("date")
	if dateStr == "" {
		c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse("date query parameter is required"))
		return
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse("invalid date format, use YYYY-MM-DD"))
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "50"))
	params := pagination.NewParams(page, pageSize)

	result, err := h.listAvailableMealsUC.Execute(c.Request.Context(), date, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederrors.NewErrorResponse(err.Error()))
		return
	}

	response := PaginatedMealResponse{
		Total:      result.Total,
		Page:       result.Page,
		PageSize:   result.PageSize,
		TotalPages: result.TotalPages,
	}
	response.Data = make([]MealResponse, len(result.Data))
	for i := range result.Data {
		response.Data[i] = ToMealResponse(&result.Data[i])
	}

	c.JSON(http.StatusOK, response)
}
