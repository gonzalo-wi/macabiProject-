package mealhttp

import (
	"net/http"
	"strconv"
	"time"

	mealusecases "macabi-back/internal/meal/application/usecases"
	sharederrors "macabi-back/internal/shared/errors"
	"macabi-back/internal/shared/pagination"
	userhttp "macabi-back/internal/user/infrastructure/http"

	"github.com/gin-gonic/gin"
)

type BookingHandler struct {
	bookMealUC        *mealusecases.BookMeal
	cancelBookingUC   *mealusecases.CancelBooking
	listMyBookingsUC  *mealusecases.ListMyBookings
	getDailySummaryUC *mealusecases.GetDailySummary
}

func NewBookingHandler(bookMealUC *mealusecases.BookMeal, cancelBookingUC *mealusecases.CancelBooking, listMyBookingsUC *mealusecases.ListMyBookings, getDailySummaryUC *mealusecases.GetDailySummary) *BookingHandler {
	return &BookingHandler{
		bookMealUC:        bookMealUC,
		cancelBookingUC:   cancelBookingUC,
		listMyBookingsUC:  listMyBookingsUC,
		getDailySummaryUC: getDailySummaryUC,
	}
}

func (h *BookingHandler) Book(c *gin.Context) {
	userID := c.GetString(userhttp.AuthUserIDKey)
	var req BookMealRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse(err.Error()))
		return
	}
	booking, err := h.bookMealUC.Execute(c.Request.Context(), mealusecases.BookMealInput{
		UserID:          userID,
		MealID:          req.MealID,
		GarnishOptionID: req.GarnishOptionID,
	})
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusCreated, ToBookingResponse(booking))
}

func (h *BookingHandler) Cancel(c *gin.Context) {
	bookingID := c.Param("id")
	userID := c.GetString(userhttp.AuthUserIDKey)
	err := h.cancelBookingUC.Execute(c.Request.Context(), mealusecases.CancelBookingInput{
		BookingID: bookingID,
		UserID:    userID,
	})
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "reserva cancelada correctamente"})
}

func (h *BookingHandler) ListMine(c *gin.Context) {
	userID := c.GetString(userhttp.AuthUserIDKey)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	params := pagination.NewParams(page, pageSize)
	result, err := h.listMyBookingsUC.Execute(c.Request.Context(), userID, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederrors.NewErrorResponse(err.Error()))
		return
	}
	response := PaginatedBookingResponse{
		Total:      result.Total,
		Page:       result.Page,
		PageSize:   result.PageSize,
		TotalPages: result.TotalPages,
	}
	response.Data = make([]BookingResponse, len(result.Data))
	for i := range result.Data {
		response.Data[i] = ToBookingResponse(&result.Data[i])
	}

	c.JSON(http.StatusOK, response)
}

func (h *BookingHandler) DailySummary(c *gin.Context) {
	dateStr := c.Query("date")
	var date time.Time
	var err error
	if dateStr == "" {
		date = time.Now().UTC()
	} else {
		date, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse("formato de fecha inválido, usar YYYY-MM-DD"))
			return
		}
	}

	projectID := c.Query("project_id")
	summary, err := h.getDailySummaryUC.Execute(c.Request.Context(), mealusecases.GetDailySummaryInput{
		Date:      date,
		ProjectID: projectID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, ToDailySummaryResponse(summary))
}
