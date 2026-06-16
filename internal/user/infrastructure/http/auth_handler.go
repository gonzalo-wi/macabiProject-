package userhttp

import (
	"net/http"

	"github.com/gin-gonic/gin"

	sharederrors "macabi-back/internal/shared/errors"
	userusecases "macabi-back/internal/user/application/usecases"
	userdto "macabi-back/internal/user/infrastructure/http/dto"
)

const forgotPasswordSuccessMessage = "Si el email está registrado, te enviamos un enlace para restablecer la contraseña."

type AuthHandler struct {
	loginUC           *userusecases.Login
	acceptInvUC       *userusecases.AcceptInvitation
	requestPasswordUC *userusecases.RequestPasswordReset
	resetPasswordUC   *userusecases.ResetPassword
}

func NewAuthHandler(
	loginUC *userusecases.Login,
	acceptInvUC *userusecases.AcceptInvitation,
	requestPasswordUC *userusecases.RequestPasswordReset,
	resetPasswordUC *userusecases.ResetPassword,
) *AuthHandler {
	return &AuthHandler{
		loginUC:           loginUC,
		acceptInvUC:       acceptInvUC,
		requestPasswordUC: requestPasswordUC,
		resetPasswordUC:   resetPasswordUC,
	}
}

func (h *AuthHandler) RegisterDisabled(c *gin.Context) {
	c.JSON(http.StatusForbidden, sharederrors.NewErrorResponse(
		"El registro público está deshabilitado. Solicitá una invitación a un administrador.",
	))
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req userdto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse(err.Error()))
		return
	}
	output, err := h.loginUC.Execute(c.Request.Context(), userusecases.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, userdto.LoginResponse{
		Token: output.Token,
		User:  userdto.ToUserResponse(output.User),
	})
}

func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req userdto.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse(err.Error()))
		return
	}
	if err := h.requestPasswordUC.Execute(c.Request.Context(), req.Email); err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": forgotPasswordSuccessMessage})
}

func (h *AuthHandler) ConfirmPasswordReset(c *gin.Context) {
	var req userdto.ConfirmPasswordResetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse(err.Error()))
		return
	}
	if err := h.resetPasswordUC.Execute(c.Request.Context(), req.Token, req.NewPassword); err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Contraseña actualizada. Ya podés iniciar sesión."})
}

func (h *AuthHandler) AcceptInvitation(c *gin.Context) {
	var req userdto.AcceptInvitationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse(err.Error()))
		return
	}
	if err := h.acceptInvUC.Execute(c.Request.Context(), req.Token, req.Password); err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Cuenta creada. Ya podés iniciar sesión."})
}
