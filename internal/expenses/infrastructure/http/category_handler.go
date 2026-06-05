package expenseshttp

import (
	"net/http"

	"github.com/gin-gonic/gin"

	expensesdto "macabi-back/internal/expenses/infrastructure/http/dto"
	sharederrors "macabi-back/internal/shared/errors"
)

func (h *Handler) CreateCategory(c *gin.Context) {
	var req expensesdto.CreateExpenseCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse(err.Error()))
		return
	}
	cat, err := h.createCategory.Execute(c.Request.Context(), req.Name)
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusCreated, expensesdto.ToExpenseCategoryResponse(*cat))
}

func (h *Handler) ListCategories(c *gin.Context) {
	cats, err := h.listCategories.Execute(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, expensesdto.ToExpenseCategoryListResponse(cats))
}

func (h *Handler) DeleteCategory(c *gin.Context) {
	if err := h.deleteCategory.Execute(c.Request.Context(), c.Param("id")); err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.Status(http.StatusNoContent)
}
