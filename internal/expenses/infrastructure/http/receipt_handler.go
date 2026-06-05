package expenseshttp

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	expensesusecases "macabi-back/internal/expenses/application/usecases"
	expensesdomain "macabi-back/internal/expenses/domain"
	expensesdto "macabi-back/internal/expenses/infrastructure/http/dto"
	sharederrors "macabi-back/internal/shared/errors"
	userhttp "macabi-back/internal/user/infrastructure/http"
)

func (h *Handler) ReceiptUpload(c *gin.Context) {
	var body expensesdto.ReceiptUploadBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse(err.Error()))
		return
	}
	res, err := h.recUpload.Execute(c.Request.Context(), expensesusecases.ReceiptUploadInput{
		ExpenseID:   c.Param("id"),
		ActorID:     c.GetString(userhttp.AuthUserIDKey),
		UserRole:    c.GetString(userhttp.AuthRoleKey),
		ContentType: body.ContentType,
	})
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"upload_url":     res.UploadURL,
		"path":           res.Path,
		"upload_headers": res.UploadHeaders,
	})
}

func (h *Handler) ReceiptUploadMultipart(c *gin.Context) {
	const maxBytes = expensesusecases.MaxReceiptBytes
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
		c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse(expensesdomain.ErrReceiptTooLarge.Error()))
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
	res, err := h.recUploadFile.Execute(c.Request.Context(), expensesusecases.ReceiptUploadFileInput{
		ExpenseID:   c.Param("id"),
		ActorID:     c.GetString(userhttp.AuthUserIDKey),
		UserRole:    c.GetString(userhttp.AuthRoleKey),
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

func (h *Handler) ReceiptView(c *gin.Context) {
	res, err := h.recView.Execute(c.Request.Context(), expensesusecases.ReceiptDownloadInput{
		ExpenseID: c.Param("id"),
		UserID:    c.GetString(userhttp.AuthUserIDKey),
		UserRole:  c.GetString(userhttp.AuthRoleKey),
	})
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, gin.H{"view_url": res.ViewURL})
}

func (h *Handler) RemoveReceipt(c *gin.Context) {
	if err := h.removeReceipt.Execute(c.Request.Context(), expensesusecases.RemoveReceiptInput{
		ExpenseID: c.Param("id"),
		ActorID:   c.GetString(userhttp.AuthUserIDKey),
		UserRole:  c.GetString(userhttp.AuthRoleKey),
	}); err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.Status(http.StatusNoContent)
}
