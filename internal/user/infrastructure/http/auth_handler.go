package userhttp

import (
	"errors"
	"net/http"

	sharederrors "macabi-back/internal/shared/errors"
	userusecases "macabi-back/internal/user/application/usecases"
	userdomain "macabi-back/internal/user/domain"

	"github.com/gin-gonic/gin"
)

const forgotPasswordSuccessMessage = "Si el email está registrado, te enviamos un enlace para restablecer la contraseña."

type AuthHandler struct {
	registerUC        *userusecases.RegisterUser
	loginUC           *userusecases.Login
	requestPasswordUC *userusecases.RequestPasswordReset
	resetPasswordUC   *userusecases.ResetPassword
}

func NewAuthHandler(
	registerUC *userusecases.RegisterUser,
	loginUC *userusecases.Login,
	requestPasswordUC *userusecases.RequestPasswordReset,
	resetPasswordUC *userusecases.ResetPassword,
) *AuthHandler {
	return &AuthHandler{
		registerUC:        registerUC,
		loginUC:           loginUC,
		requestPasswordUC: requestPasswordUC,
		resetPasswordUC:   resetPasswordUC,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse(err.Error()))
		return
	}

	user, err := h.registerUC.Execute(c.Request.Context(), userusecases.RegisterInput{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, ToUserResponse(user))
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
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

	c.JSON(http.StatusOK, LoginResponse{
		Token: output.Token,
		User:  ToUserResponse(output.User),
	})
}

func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse(err.Error()))
		return
	}

	err := h.requestPasswordUC.Execute(c.Request.Context(), req.Email)
	if err != nil {
		if errors.Is(err, userdomain.ErrInvalidEmail) {
			c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse(err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, sharederrors.NewErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": forgotPasswordSuccessMessage})
}

func (h *AuthHandler) ConfirmPasswordReset(c *gin.Context) {
	var req ConfirmPasswordResetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse(err.Error()))
		return
	}

	err := h.resetPasswordUC.Execute(c.Request.Context(), req.Token, req.NewPassword)
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Contraseña actualizada. Ya podés iniciar sesión."})
}
