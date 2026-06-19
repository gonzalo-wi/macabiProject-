package newshttp

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	newsusecases "macabi-back/internal/news/application/usecases"
	newsdomain "macabi-back/internal/news/domain"
	newsdto "macabi-back/internal/news/infrastructure/http/dto"
	sharederrors "macabi-back/internal/shared/errors"
	"macabi-back/internal/shared/pagination"
	userhttp "macabi-back/internal/user/infrastructure/http"
)

func (h *Handler) Create(c *gin.Context) {
	if isMultipartRequest(c) {
		h.createMultipart(c)
		return
	}
	var body newsdto.CreateNewsBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse(err.Error()))
		return
	}
	n, err := h.create.Execute(c.Request.Context(), newsusecases.CreateNewsInput{
		Title:      body.Title,
		Body:       body.Body,
		AuthorID:   c.GetString(userhttp.AuthUserIDKey),
		Publish:    body.Publish,
		ProjectIDs: body.ProjectIDs,
	})
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusCreated, h.toResp(c.Request.Context(), *n))
}

func (h *Handler) createMultipart(c *gin.Context) {
	const maxBytes = newsusecases.MaxImageBytes
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBytes+64*1024)
	if err := c.Request.ParseMultipartForm(maxBytes + 64*1024); err != nil {
		c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse("no se pudo leer el formulario"))
		return
	}
	base := newsusecases.CreateNewsInput{
		Title:      strings.TrimSpace(c.PostForm("title")),
		Body:       strings.TrimSpace(c.PostForm("body")),
		AuthorID:   c.GetString(userhttp.AuthUserIDKey),
		Publish:    parseBool(c.PostForm("publish")),
		ProjectIDs: c.PostFormArray("project_ids"),
	}

	file, err := c.FormFile("file")
	if err != nil {
		n, err := h.create.Execute(c.Request.Context(), base)
		if err != nil {
			c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
			return
		}
		c.JSON(http.StatusCreated, h.toResp(c.Request.Context(), *n))
		return
	}
	if file.Size > maxBytes {
		c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse(newsdomain.ErrImageTooLarge.Error()))
		return
	}
	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse("no se pudo leer la imagen"))
		return
	}
	defer f.Close()
	contentType := strings.TrimSpace(file.Header.Get("Content-Type"))
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	n, err := h.createWithImage.Execute(c.Request.Context(), newsusecases.CreateNewsWithImageInput{
		CreateNewsInput:  base,
		ImageContentType: contentType,
		ImageSize:        file.Size,
		ImageBody:        f,
	})
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusCreated, h.toResp(c.Request.Context(), *n))
}

func (h *Handler) Get(c *gin.Context) {
	n, err := h.getNews.Execute(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, h.toResp(c.Request.Context(), *n))
}

func (h *Handler) ListPublished(c *gin.Context) {
	params := pagination.ParseParams(c.Query("page"), c.Query("page_size"))
	res, err := h.listPublished.Execute(c.Request.Context(), c.GetString(userhttp.AuthUserIDKey), params)
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, h.toListResp(c.Request.Context(), res))
}

func (h *Handler) ListAll(c *gin.Context) {
	params := pagination.ParseParams(c.Query("page"), c.Query("page_size"))
	res, err := h.listAll.Execute(c.Request.Context(), params)
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, h.toListResp(c.Request.Context(), res))
}

// Latest devuelve la última noticia publicada, o null si todavía no hay ninguna.
func (h *Handler) Latest(c *gin.Context) {
	n, err := h.getLatest.Execute(c.Request.Context(), c.GetString(userhttp.AuthUserIDKey))
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	if n == nil {
		c.JSON(http.StatusOK, nil)
		return
	}
	c.JSON(http.StatusOK, h.toResp(c.Request.Context(), *n))
}

func (h *Handler) Patch(c *gin.Context) {
	var body newsdto.PatchNewsBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse(err.Error()))
		return
	}
	n, err := h.upd.Execute(c.Request.Context(), newsusecases.UpdateNewsInput{
		ID:         c.Param("id"),
		Title:      body.Title,
		Body:       body.Body,
		Publish:    body.Publish,
		Renotify:   body.Renotify,
		ProjectIDs: body.ProjectIDs,
	})
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, h.toResp(c.Request.Context(), *n))
}

func (h *Handler) Delete(c *gin.Context) {
	if err := h.deleteUC.Execute(c.Request.Context(), c.Param("id")); err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) ImageUpload(c *gin.Context) {
	const maxBytes = newsusecases.MaxImageBytes
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBytes+512)
	if err := c.Request.ParseMultipartForm(maxBytes + 512); err != nil {
		c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse("no se pudo leer el formulario"))
		return
	}
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse("archivo requerido"))
		return
	}
	if file.Size > maxBytes {
		c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse(newsdomain.ErrImageTooLarge.Error()))
		return
	}
	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse("no se pudo leer el archivo"))
		return
	}
	defer f.Close()
	contentType := strings.TrimSpace(file.Header.Get("Content-Type"))
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	res, err := h.imageUpload.Execute(c.Request.Context(), newsusecases.ImageUploadFileInput{
		NewsID:      c.Param("id"),
		ContentType: contentType,
		Size:        file.Size,
		Body:        f,
	})
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, gin.H{"path": res.Path})
}

func (h *Handler) ImageView(c *gin.Context) {
	url, err := h.imageView.Execute(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, gin.H{"view_url": url})
}

func (h *Handler) RemoveImage(c *gin.Context) {
	if err := h.removeImage.Execute(c.Request.Context(), c.Param("id")); err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.Status(http.StatusNoContent)
}

func parseBool(s string) bool {
	switch strings.TrimSpace(strings.ToLower(s)) {
	case "true", "1", "yes", "on":
		return true
	default:
		return false
	}
}
