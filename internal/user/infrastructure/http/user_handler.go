package userhttp

import (
	"net/http"
	"strconv"

	sharederrors "macabi-back/internal/shared/errors"
	"macabi-back/internal/shared/pagination"
	userusecases "macabi-back/internal/user/application/usecases"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	getCurrentUserUC *userusecases.GetCurrentUser
	changeRoleUC     *userusecases.ChangeRole
	listUsersUC      *userusecases.ListUsers
	setUserStatusUC  *userusecases.SetUserStatus
	updateUserUC     *userusecases.UpdateUser
	changePasswordUC *userusecases.ChangePassword
}

func NewUserHandler(
	getCurrentUserUC *userusecases.GetCurrentUser,
	changeRoleUC *userusecases.ChangeRole,
	listUsersUC *userusecases.ListUsers,
	setUserStatusUC *userusecases.SetUserStatus,
	updateUserUC *userusecases.UpdateUser,
	changePasswordUC *userusecases.ChangePassword,
) *UserHandler {
	return &UserHandler{
		getCurrentUserUC: getCurrentUserUC,
		changeRoleUC:     changeRoleUC,
		listUsersUC:      listUsersUC,
		setUserStatusUC:  setUserStatusUC,
		updateUserUC:     updateUserUC,
		changePasswordUC: changePasswordUC,
	}
}

func (h *UserHandler) Me(c *gin.Context) {
	userID := c.GetString(AuthUserIDKey)
	user, err := h.getCurrentUserUC.Execute(c.Request.Context(), userID)
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, ToUserResponse(user))
}

func (h *UserHandler) ChangeRole(c *gin.Context) {
	targetUserID := c.Param("id")
	changedByID := c.GetString(AuthUserIDKey)
	var req ChangeRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse(err.Error()))
		return
	}

	err := h.changeRoleUC.Execute(c.Request.Context(), userusecases.ChangeRoleInput{
		TargetUserID: targetUserID,
		NewRole:      req.Role,
		ChangedByID:  changedByID,
	})
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "rol actualizado correctamente"})
}

func (h *UserHandler) ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	params := pagination.NewParams(page, pageSize)

	result, err := h.listUsersUC.Execute(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederrors.NewErrorResponse(err.Error()))
		return
	}

	response := PaginatedUserResponse{
		Total:      result.Total,
		Page:       result.Page,
		PageSize:   result.PageSize,
		TotalPages: result.TotalPages,
	}
	response.Data = make([]UserResponse, len(result.Data))
	for i := range result.Data {
		response.Data[i] = ToUserResponse(&result.Data[i])
	}

	c.JSON(http.StatusOK, response)
}

func (h *UserHandler) SetStatus(c *gin.Context) {
	targetUserID := c.Param("id")
	var req SetUserStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse(err.Error()))
		return
	}

	err := h.setUserStatusUC.Execute(c.Request.Context(), userusecases.SetUserStatusInput{
		TargetUserID: targetUserID,
		Active:       req.Active,
	})
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "estado actualizado correctamente"})
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	targetUserID := c.Param("id")
	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse(err.Error()))
		return
	}

	user, err := h.updateUserUC.Execute(c.Request.Context(), userusecases.UpdateUserInput{
		TargetUserID: targetUserID,
		Name:         req.Name,
		Email:        req.Email,
	})
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, ToUserResponse(user))
}

func (h *UserHandler) ChangePassword(c *gin.Context) {
	userID := c.GetString(AuthUserIDKey)
	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse(err.Error()))
		return
	}

	err := h.changePasswordUC.Execute(c.Request.Context(), userusecases.ChangePasswordInput{
		UserID:          userID,
		CurrentPassword: req.CurrentPassword,
		NewPassword:     req.NewPassword,
	})
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "contraseña actualizada correctamente"})
}
