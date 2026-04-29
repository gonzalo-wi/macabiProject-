package userhttp

import (
	"net/http"

	sharederrors "macabi-back/internal/shared/errors"
	userusecases "macabi-back/internal/user/application/usecases"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	loginUC     *userusecases.Login
	acceptInvUC *userusecases.AcceptInvitation
}

func NewAuthHandler(
	loginUC *userusecases.Login,
	acceptInvUC *userusecases.AcceptInvitation,
) *AuthHandler {
	return &AuthHandler{
		loginUC:     loginUC,
		acceptInvUC: acceptInvUC,
	}
}

func (h *AuthHandler) RegisterDisabled(c *gin.Context) {
	c.JSON(http.StatusForbidden, sharederrors.NewErrorResponse(
		"El registro público está deshabilitado. Solicitá una invitación a un administrador.",
	))
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

func (h *AuthHandler) AcceptInvitation(c *gin.Context) {
	var req AcceptInvitationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse(err.Error()))
		return
	}

	err := h.acceptInvUC.Execute(c.Request.Context(), req.Token, req.Password)
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Cuenta creada. Ya podés iniciar sesión."})
}
