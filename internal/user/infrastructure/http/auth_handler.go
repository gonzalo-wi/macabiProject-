package userhttp

import (
	"net/http"

	userusecases "macabi-back/internal/user/application/usecases"
	sharederrors "macabi-back/internal/shared/errors"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	registerUC *userusecases.RegisterUser
	loginUC    *userusecases.Login
}

func NewAuthHandler(registerUC *userusecases.RegisterUser, loginUC *userusecases.Login) *AuthHandler {
	return &AuthHandler{registerUC: registerUC, loginUC: loginUC}
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
