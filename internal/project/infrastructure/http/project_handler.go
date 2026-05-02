package projecthttp

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	projectusecases "macabi-back/internal/project/application/usecases"
	sharederrors "macabi-back/internal/shared/errors"
	"macabi-back/internal/shared/pagination"
)

type ProjectHandler struct {
	createUC *projectusecases.CreateProject
	listUC   *projectusecases.ListProjects
	getUC    *projectusecases.GetProject
	updateUC *projectusecases.UpdateProject
	deleteUC *projectusecases.DeleteProject
}

func NewProjectHandler(
	createUC *projectusecases.CreateProject,
	listUC *projectusecases.ListProjects,
	getUC *projectusecases.GetProject,
	updateUC *projectusecases.UpdateProject,
	deleteUC *projectusecases.DeleteProject,
) *ProjectHandler {
	return &ProjectHandler{
		createUC: createUC,
		listUC:   listUC,
		getUC:    getUC,
		updateUC: updateUC,
		deleteUC: deleteUC,
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
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	params := pagination.NewParams(page, pageSize)
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
