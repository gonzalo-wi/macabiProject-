package projecthttp

import (
	"net/http"

	"github.com/gin-gonic/gin"

	projectusecases "macabi-back/internal/project/application/usecases"
	projectdto "macabi-back/internal/project/infrastructure/http/dto"
	sharederrors "macabi-back/internal/shared/errors"
)

func (h *ProjectHandler) ListMembers(c *gin.Context) {
	list, err := h.listMembersUC.Execute(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	out := make([]projectdto.ProjectMemberResponse, len(list))
	for i := range list {
		out[i] = projectdto.ToMemberResponse(&list[i])
	}
	c.JSON(http.StatusOK, gin.H{"data": out})
}

func (h *ProjectHandler) AddMember(c *gin.Context) {
	var req projectdto.AddProjectMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse(err.Error()))
		return
	}
	m, err := h.addMemberUC.Execute(c.Request.Context(), projectusecases.AddProjectMemberInput{
		ProjectID: c.Param("id"),
		UserID:    req.UserID,
		Role:      req.Role,
	})
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusCreated, projectdto.ToMemberResponse(m))
}

func (h *ProjectHandler) RemoveMember(c *gin.Context) {
	if err := h.removeMemberUC.Execute(c.Request.Context(), c.Param("id"), c.Param("userId")); err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.Status(http.StatusNoContent)
}
