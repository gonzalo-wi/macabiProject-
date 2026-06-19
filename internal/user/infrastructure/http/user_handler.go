package userhttp

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	sharederrors "macabi-back/internal/shared/errors"
	"macabi-back/internal/shared/pagination"
	userports "macabi-back/internal/user/application/ports"
	userusecases "macabi-back/internal/user/application/usecases"
	userdomain "macabi-back/internal/user/domain"
	userdto "macabi-back/internal/user/infrastructure/http/dto"
)

type UserHandler struct {
	getCurrentUserUC         *userusecases.GetCurrentUser
	changeRoleUC             *userusecases.ChangeRole
	listUsersUC              *userusecases.ListUsers
	setUserStatusUC          *userusecases.SetUserStatus
	updateUserUC             *userusecases.UpdateUser
	changePasswordUC         *userusecases.ChangePassword
	createUserByAdminUC        *userusecases.CreateUserByAdmin
	deleteUserByAdminUC        *userusecases.DeleteUserByAdmin
	createUserInvitation     *userusecases.CreateUserInvitation
	listPendingInvitationsUC *userusecases.ListPendingInvitations
	resendUserInvitationUC   *userusecases.ResendUserInvitation
	revokeUserInvitationUC   *userusecases.RevokeUserInvitation
}

func NewUserHandler(
	getCurrentUserUC *userusecases.GetCurrentUser,
	changeRoleUC *userusecases.ChangeRole,
	listUsersUC *userusecases.ListUsers,
	setUserStatusUC *userusecases.SetUserStatus,
	updateUserUC *userusecases.UpdateUser,
	changePasswordUC *userusecases.ChangePassword,
	createUserByAdminUC *userusecases.CreateUserByAdmin,
	deleteUserByAdminUC *userusecases.DeleteUserByAdmin,
	createUserInvitation *userusecases.CreateUserInvitation,
	listPendingInvitationsUC *userusecases.ListPendingInvitations,
	resendUserInvitationUC *userusecases.ResendUserInvitation,
	revokeUserInvitationUC *userusecases.RevokeUserInvitation,
) *UserHandler {
	return &UserHandler{
		getCurrentUserUC:         getCurrentUserUC,
		changeRoleUC:             changeRoleUC,
		listUsersUC:              listUsersUC,
		setUserStatusUC:          setUserStatusUC,
		updateUserUC:             updateUserUC,
		changePasswordUC:         changePasswordUC,
		createUserByAdminUC:        createUserByAdminUC,
		deleteUserByAdminUC:        deleteUserByAdminUC,
		createUserInvitation:     createUserInvitation,
		listPendingInvitationsUC: listPendingInvitationsUC,
		resendUserInvitationUC:   resendUserInvitationUC,
		revokeUserInvitationUC:   revokeUserInvitationUC,
	}
}

func (h *UserHandler) Me(c *gin.Context) {
	user, err := h.getCurrentUserUC.Execute(c.Request.Context(), c.GetString(AuthUserIDKey))
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, userdto.ToUserResponse(user))
}

func (h *UserHandler) ChangeRole(c *gin.Context) {
	var req userdto.ChangeRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse(err.Error()))
		return
	}
	if err := h.changeRoleUC.Execute(c.Request.Context(), userusecases.ChangeRoleInput{
		TargetUserID: c.Param("id"),
		NewRole:      req.Role,
		ChangedByID:  c.GetString(AuthUserIDKey),
	}); err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "rol actualizado correctamente"})
}

func (h *UserHandler) ListUsers(c *gin.Context) {
	params := pagination.ParseParams(c.Query("page"), c.Query("page_size"))
	filter := userports.UserListFilter{Query: strings.TrimSpace(c.Query("q"))}
	result, err := h.listUsersUC.Execute(c.Request.Context(), filter, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederrors.NewErrorResponse(err.Error()))
		return
	}
	response := userdto.PaginatedUserResponse{
		Total:      result.Total,
		Page:       result.Page,
		PageSize:   result.PageSize,
		TotalPages: result.TotalPages,
		Data:       make([]userdto.UserResponse, len(result.Data)),
	}
	for i := range result.Data {
		row := result.Data[i]
		response.Data[i] = userdto.ToUserResponseWithInvitation(
			&row.User,
			row.PendingInvitationID,
			row.InvitationStatus,
		)
	}
	c.JSON(http.StatusOK, response)
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req userdto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse(err.Error()))
		return
	}
	user, err := h.createUserByAdminUC.Execute(c.Request.Context(), userusecases.CreateUserByAdminInput{
		Name:          req.Name,
		Email:         req.Email,
		RequestedRole: req.Role,
		InviterRole:   userdomain.Role(c.GetString(AuthRoleKey)),
	})
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusCreated, userdto.ToUserResponseWithInvitation(user, "", userdomain.InvitationStatus(user, false)))
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	if err := h.deleteUserByAdminUC.Execute(c.Request.Context(), userusecases.DeleteUserByAdminInput{
		TargetUserID: c.Param("id"),
		ActorUserID:  c.GetString(AuthUserIDKey),
	}); err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "usuario eliminado correctamente"})
}

func (h *UserHandler) SetStatus(c *gin.Context) {
	var req userdto.SetUserStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse(err.Error()))
		return
	}
	if err := h.setUserStatusUC.Execute(c.Request.Context(), userusecases.SetUserStatusInput{
		TargetUserID: c.Param("id"),
		Active:       req.Active,
	}); err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "estado actualizado correctamente"})
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	var req userdto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse(err.Error()))
		return
	}
	user, err := h.updateUserUC.Execute(c.Request.Context(), userusecases.UpdateUserInput{
		TargetUserID: c.Param("id"),
		Name:         req.Name,
		Email:        req.Email,
	})
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, userdto.ToUserResponse(user))
}

func (h *UserHandler) ChangePassword(c *gin.Context) {
	var req userdto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse(err.Error()))
		return
	}
	if err := h.changePasswordUC.Execute(c.Request.Context(), userusecases.ChangePasswordInput{
		UserID:          c.GetString(AuthUserIDKey),
		CurrentPassword: req.CurrentPassword,
		NewPassword:     req.NewPassword,
	}); err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "contraseña actualizada correctamente"})
}
