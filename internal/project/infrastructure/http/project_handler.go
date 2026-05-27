package projecthttp

import (
	"net/http"

	"github.com/gin-gonic/gin"

	projectusecases "macabi-back/internal/project/application/usecases"
	sharederrors "macabi-back/internal/shared/errors"
	"macabi-back/internal/shared/pagination"
)

type ProjectHandler struct {
	createUC       *projectusecases.CreateProject
	listUC         *projectusecases.ListProjects
	getUC          *projectusecases.GetProject
	updateUC       *projectusecases.UpdateProject
	deleteUC       *projectusecases.DeleteProject
	addMemberUC    *projectusecases.AddProjectMember
	removeMemberUC *projectusecases.RemoveProjectMember
	listMembersUC  *projectusecases.ListProjectMembers
}

func NewProjectHandler(
	createUC *projectusecases.CreateProject,
	listUC *projectusecases.ListProjects,
	getUC *projectusecases.GetProject,
	updateUC *projectusecases.UpdateProject,
	deleteUC *projectusecases.DeleteProject,
	addMemberUC *projectusecases.AddProjectMember,
	removeMemberUC *projectusecases.RemoveProjectMember,
	listMembersUC *projectusecases.ListProjectMembers,
) *ProjectHandler {
	return &ProjectHandler{
		createUC:       createUC,
		listUC:         listUC,
		getUC:          getUC,
		updateUC:       updateUC,
		deleteUC:       deleteUC,
		addMemberUC:    addMemberUC,
		removeMemberUC: removeMemberUC,
		listMembersUC:  listMembersUC,
	}
}

func (h *ProjectHandler) Create(c *gin.Context) {
	var req CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse(err.Error()))
		return
	}
	p, err := h.createUC.Execute(c.Request.Context(), req.toInput())
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusCreated, toProjectResponse(p))
}

func (h *ProjectHandler) List(c *gin.Context) {
	params := pagination.ParseParams(c.Query("page"), c.Query("page_size"))
	result, err := h.listUC.Execute(c.Request.Context(), params)
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, toProjectListResponse(result))
}

func (h *ProjectHandler) Get(c *gin.Context) {
	p, err := h.getUC.Execute(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, toProjectResponse(p))
}

func (h *ProjectHandler) Update(c *gin.Context) {
	var req UpdateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse(err.Error()))
		return
	}
	p, err := h.updateUC.Execute(c.Request.Context(), req.toInput(c.Param("id")))
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, toProjectResponse(p))
}

func (h *ProjectHandler) Delete(c *gin.Context) {
	if err := h.deleteUC.Execute(c.Request.Context(), c.Param("id")); err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *ProjectHandler) ListMembers(c *gin.Context) {
	list, err := h.listMembersUC.Execute(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	out := make([]ProjectMemberResponse, len(list))
	for i := range list {
		out[i] = toMemberResponse(&list[i])
	}
	c.JSON(http.StatusOK, gin.H{"data": out})
}

func (h *ProjectHandler) AddMember(c *gin.Context) {
	var req AddProjectMemberRequest
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
	c.JSON(http.StatusCreated, toMemberResponse(m))
}

func (h *ProjectHandler) RemoveMember(c *gin.Context) {
	if err := h.removeMemberUC.Execute(c.Request.Context(), c.Param("id"), c.Param("userId")); err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.Status(http.StatusNoContent)
}
