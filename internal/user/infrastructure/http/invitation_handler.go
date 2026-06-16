package userhttp

import (
	"net/http"

	"github.com/gin-gonic/gin"

	sharederrors "macabi-back/internal/shared/errors"
	userusecases "macabi-back/internal/user/application/usecases"
	userdomain "macabi-back/internal/user/domain"
	userdto "macabi-back/internal/user/infrastructure/http/dto"
)

type InvitationHandler struct {
	createUC  *userusecases.CreateUserInvitation
	listUC    *userusecases.ListPendingInvitations
	resendUC  *userusecases.ResendUserInvitation
	revokeUC  *userusecases.RevokeUserInvitation
}

func NewInvitationHandler(
	createUC *userusecases.CreateUserInvitation,
	listUC *userusecases.ListPendingInvitations,
	resendUC *userusecases.ResendUserInvitation,
	revokeUC *userusecases.RevokeUserInvitation,
) *InvitationHandler {
	return &InvitationHandler{
		createUC: createUC,
		listUC:   listUC,
		resendUC: resendUC,
		revokeUC: revokeUC,
	}
}

func (h *InvitationHandler) Create(c *gin.Context) {
	var req userdto.CreateUserInvitationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse(err.Error()))
		return
	}
	if err := h.createUC.Execute(c.Request.Context(), userusecases.CreateUserInvitationInput{
		Email:         req.Email,
		Name:          req.Name,
		RequestedRole: req.Role,
		InviterID:     c.GetString(AuthUserIDKey),
		InviterRole:   userdomain.Role(c.GetString(AuthRoleKey)),
	}); err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Invitación enviada correctamente."})
}

func (h *InvitationHandler) ListPending(c *gin.Context) {
	list, err := h.listUC.Execute(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederrors.NewErrorResponse(err.Error()))
		return
	}
	out := userdto.ListPendingInvitationsResponse{
		Data: make([]userdto.PendingInvitationResponse, len(list)),
	}
	for i := range list {
		out.Data[i] = userdto.ToPendingInvitationResponse(list[i])
	}
	c.JSON(http.StatusOK, out)
}

func (h *InvitationHandler) Resend(c *gin.Context) {
	if err := h.resendUC.Execute(c.Request.Context(), c.Param("id")); err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Invitación reenviada correctamente."})
}

func (h *InvitationHandler) Revoke(c *gin.Context) {
	if err := h.revokeUC.Execute(c.Request.Context(), c.Param("id")); err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Invitación revocada correctamente."})
}
