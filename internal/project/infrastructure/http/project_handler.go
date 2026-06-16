package projecthttp

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	projectdto "macabi-back/internal/project/infrastructure/http/dto"
	projectports "macabi-back/internal/project/application/ports"
	sharederrors "macabi-back/internal/shared/errors"
	"macabi-back/internal/shared/pagination"
)

func (h *ProjectHandler) Create(c *gin.Context) {
	var req projectdto.CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse(err.Error()))
		return
	}
	p, err := h.createUC.Execute(c.Request.Context(), req.ToInput())
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusCreated, projectdto.ToProjectResponse(p))
}

func (h *ProjectHandler) List(c *gin.Context) {
	params := pagination.ParseParams(c.Query("page"), c.Query("page_size"))
	filter := projectports.ProjectListFilter{Query: strings.TrimSpace(c.Query("q"))}
	result, err := h.listUC.Execute(c.Request.Context(), filter, params)
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, projectdto.ToProjectListResponse(result))
}

func (h *ProjectHandler) Get(c *gin.Context) {
	p, err := h.getUC.Execute(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, projectdto.ToProjectResponse(p))
}

func (h *ProjectHandler) Update(c *gin.Context) {
	var req projectdto.UpdateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse(err.Error()))
		return
	}
	p, err := h.updateUC.Execute(c.Request.Context(), req.ToInput(c.Param("id")))
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, projectdto.ToProjectResponse(p))
}

func (h *ProjectHandler) Delete(c *gin.Context) {
	if err := h.deleteUC.Execute(c.Request.Context(), c.Param("id")); err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.Status(http.StatusNoContent)
}
