package expenseshttp

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	expensesusecases "macabi-back/internal/expenses/application/usecases"
	expensesdto "macabi-back/internal/expenses/infrastructure/http/dto"
	sharederrors "macabi-back/internal/shared/errors"
	userhttp "macabi-back/internal/user/infrastructure/http"
)

func (h *Handler) Analytics(c *gin.Context) {
	var fromPtr, toPtr *time.Time
	if f := strings.TrimSpace(c.Query("from")); f != "" {
		t0, err := time.Parse("2006-01-02", f)
		if err != nil {
			c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse(err.Error()))
			return
		}
		tUTC := t0.UTC().Truncate(24 * time.Hour)
		fromPtr = &tUTC
	}
	if f := strings.TrimSpace(c.Query("to")); f != "" {
		t1, err := time.Parse("2006-01-02", f)
		if err != nil {
			c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse(err.Error()))
			return
		}
		tUTC := t1.UTC().Truncate(24 * time.Hour)
		toPtr = &tUTC
	}
	out, err := h.analytics.Execute(c.Request.Context(), expensesusecases.ExpenseAnalyticsInput{From: fromPtr, To: toPtr})
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, expensesdto.AnalyticsToResp(out))
}

func (h *Handler) SummaryByProject(c *gin.Context) {
	var fromPtr, toPtr *time.Time
	if f := strings.TrimSpace(c.Query("from")); f != "" {
		t0, err := time.Parse("2006-01-02", f)
		if err != nil {
			c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse(err.Error()))
			return
		}
		tUTC := t0.UTC().Truncate(24 * time.Hour)
		fromPtr = &tUTC
	}
	if f := strings.TrimSpace(c.Query("to")); f != "" {
		t1, err := time.Parse("2006-01-02", f)
		if err != nil {
			c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse(err.Error()))
			return
		}
		tUTC := t1.UTC().Truncate(24 * time.Hour)
		toPtr = &tUTC
	}
	out, err := h.summaryUC.Execute(c.Request.Context(), expensesusecases.SummaryInput{
		ProjectID: c.Param("id"),
		UserID:    c.GetString(userhttp.AuthUserIDKey),
		UserRole:  c.GetString(userhttp.AuthRoleKey),
		From:      fromPtr,
		To:        toPtr,
	})
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, expensesdto.SummaryToResp(out))
}
