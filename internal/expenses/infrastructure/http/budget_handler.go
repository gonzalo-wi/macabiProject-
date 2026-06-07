package expenseshttp

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"

	expensesdto "macabi-back/internal/expenses/infrastructure/http/dto"
	sharederrors "macabi-back/internal/shared/errors"
	userhttp "macabi-back/internal/user/infrastructure/http"
)

func (h *Handler) GetBudget(c *gin.Context) {
	out, err := h.getBudget.Execute(
		c.Request.Context(),
		c.Param("id"),
		c.GetString(userhttp.AuthUserIDKey),
		c.GetString(userhttp.AuthRoleKey),
	)
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, expensesdto.BudgetToResp(out))
}

func (h *Handler) SetBudget(c *gin.Context) {
	var req expensesdto.SetProjectBudgetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse(err.Error()))
		return
	}
	var amount *decimal.Decimal
	if req.MonthlyAmount != nil && strings.TrimSpace(*req.MonthlyAmount) != "" {
		d, err := decimal.NewFromString(strings.TrimSpace(*req.MonthlyAmount))
		if err != nil {
			c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse("monto inválido"))
			return
		}
		amount = &d
	}
	if err := h.setBudget.Execute(c.Request.Context(), c.Param("id"), amount); err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.Status(http.StatusNoContent)
}
