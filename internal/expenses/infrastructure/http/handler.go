package expenseshttp

import (
	"mime"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	expensesusecases "macabi-back/internal/expenses/application/usecases"
	expensesdomain "macabi-back/internal/expenses/domain"
	sharederrors "macabi-back/internal/shared/errors"
	"macabi-back/internal/shared/pagination"
	userhttp "macabi-back/internal/user/infrastructure/http"
)

type Handler struct {
	create            *expensesusecases.CreateExpense
	createWithReceipt *expensesusecases.CreateExpenseWithReceipt
	getExp    *expensesusecases.GetExpense
	listProj  *expensesusecases.ListProjectExpenses
	listMine  *expensesusecases.ListMyExpenses
	upd       *expensesusecases.UpdateExpense
	approve   *expensesusecases.ApproveExpense
	reject    *expensesusecases.RejectExpense
	deleteUC  *expensesusecases.DeleteExpense
	summaryUC *expensesusecases.ProjectExpenseSummaryUC
	recUpload     *expensesusecases.ReceiptUploadURL
	recUploadFile *expensesusecases.ReceiptUploadFile
	recView       *expensesusecases.ReceiptDownloadURL
	listNotifs    *expensesusecases.ListExpenseNotifications
	markNotifRead *expensesusecases.MarkExpenseNotificationRead
	unreadNotifs  *expensesusecases.UnreadExpenseNotificationCount
}

func NewHandler(
	create *expensesusecases.CreateExpense,
	createWithReceipt *expensesusecases.CreateExpenseWithReceipt,
	getExp *expensesusecases.GetExpense,
	listProj *expensesusecases.ListProjectExpenses,
	listMine *expensesusecases.ListMyExpenses,
	upd *expensesusecases.UpdateExpense,
	approve *expensesusecases.ApproveExpense,
	reject *expensesusecases.RejectExpense,
	deleteUC *expensesusecases.DeleteExpense,
	summaryUC *expensesusecases.ProjectExpenseSummaryUC,
	recUpload *expensesusecases.ReceiptUploadURL,
	recUploadFile *expensesusecases.ReceiptUploadFile,
	recView *expensesusecases.ReceiptDownloadURL,
	listNotifs *expensesusecases.ListExpenseNotifications,
	markNotifRead *expensesusecases.MarkExpenseNotificationRead,
	unreadNotifs *expensesusecases.UnreadExpenseNotificationCount,
) *Handler {
	return &Handler{
		create:            create,
		createWithReceipt: createWithReceipt,
		getExp:    getExp,
		listProj:  listProj,
		listMine:  listMine,
		upd:       upd,
		approve:   approve,
		reject:    reject,
		deleteUC:  deleteUC,
		summaryUC: summaryUC,
		recUpload:     recUpload,
		recUploadFile: recUploadFile,
		recView:       recView,
		listNotifs:    listNotifs,
		markNotifRead: markNotifRead,
		unreadNotifs:  unreadNotifs,
	}
}

func isMultipartRequest(c *gin.Context) bool {
	ct := strings.TrimSpace(c.GetHeader("Content-Type"))
	if ct == "" {
		return false
	}
	mediaType, _, err := mime.ParseMediaType(ct)
	if err != nil {
		return strings.Contains(strings.ToLower(ct), "multipart/form-data")
	}
	return mediaType == "multipart/form-data"
}

func (h *Handler) Create(c *gin.Context) {
	if isMultipartRequest(c) {
		h.createMultipart(c)
		return
	}

	var body CreateExpenseBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse(err.Error()))
		return
	}

	exp, err := h.createExpenseJSON(c, body)
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusCreated, expenseToResp(*exp, ""))
}

func (h *Handler) createExpenseJSON(c *gin.Context, body CreateExpenseBody) (*expensesdomain.Expense, error) {
	d, err := parseExpenseDate("2006-01-02", body.ExpenseDate)
	if err != nil {
		return nil, err
	}
	return h.create.Execute(c.Request.Context(), expensesusecases.CreateExpenseInput{
		ProjectID:   body.ProjectID,
		SubmittedBy: c.GetString(userhttp.AuthUserIDKey),
		UserRole:    c.GetString(userhttp.AuthRoleKey),
		AmountStr:   body.Amount,
		Description: body.Description,
		ExpenseDate: d,
		Currency:    body.Currency,
	})
}

func (h *Handler) createMultipart(c *gin.Context) {
	const maxBytes = expensesusecases.MaxReceiptBytes
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBytes+64*1024)

	if err := c.Request.ParseMultipartForm(maxBytes + 64*1024); err != nil {
		c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse("no se pudo leer el formulario"))
		return
	}

	projectID := strings.TrimSpace(c.PostForm("project_id"))
	amount := strings.TrimSpace(c.PostForm("amount"))
	description := strings.TrimSpace(c.PostForm("description"))
	expenseDate := strings.TrimSpace(c.PostForm("expense_date"))
	currency := strings.TrimSpace(c.PostForm("currency"))

	d, err := parseExpenseDate("2006-01-02", expenseDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse(err.Error()))
		return
	}

	base := expensesusecases.CreateExpenseInput{
		ProjectID:   projectID,
		SubmittedBy: c.GetString(userhttp.AuthUserIDKey),
		UserRole:    c.GetString(userhttp.AuthRoleKey),
		AmountStr:   amount,
		Description: description,
		ExpenseDate: d,
		Currency:    currency,
	}

	file, err := c.FormFile("file")
	if err != nil {
		exp, err := h.create.Execute(c.Request.Context(), base)
		if err != nil {
			c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
			return
		}
		c.JSON(http.StatusCreated, expenseToResp(*exp, ""))
		return
	}

	if file.Size > maxBytes {
		c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse(expensesdomain.ErrReceiptTooLarge.Error()))
		return
	}

	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse("no se pudo leer el comprobante"))
		return
	}
	defer f.Close()

	contentType := strings.TrimSpace(file.Header.Get("Content-Type"))
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	exp, err := h.createWithReceipt.Execute(c.Request.Context(), expensesusecases.CreateExpenseWithReceiptInput{
		CreateExpenseInput: base,
		ReceiptContentType: contentType,
		ReceiptSize:        file.Size,
		ReceiptBody:        f,
	})
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusCreated, expenseToResp(*exp, ""))
}

func (h *Handler) ListMine(c *gin.Context) {
	params := pagination.NewParams(parsePage(c.DefaultQuery("page", "1"), 1), parsePage(c.DefaultQuery("page_size", "20"), 20))
	res, err := h.listMine.Execute(c.Request.Context(), c.GetString(userhttp.AuthUserIDKey), params)
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, listResp(res))
}

func (h *Handler) ListByProject(c *gin.Context) {
	pid := c.Param("id")
	params := pagination.NewParams(parsePage(c.DefaultQuery("page", "1"), 1), parsePage(c.DefaultQuery("page_size", "20"), 20))

	res, err := h.listProj.Execute(c.Request.Context(), expensesusecases.ListProjectExpensesInput{
		ProjectID: pid,
		UserID:    c.GetString(userhttp.AuthUserIDKey),
		UserRole:  c.GetString(userhttp.AuthRoleKey),
		Params:    params,
	})
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, listResp(res))
}

func (h *Handler) SummaryByProject(c *gin.Context) {
	pid := c.Param("id")

	var fromPtr, toPtr *time.Time
	if f := strings.TrimSpace(c.Query("from")); f != "" {
		t0, err := time.Parse("2006-01-02", f)
		if err != nil {
			c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse(err.Error()))
			return
		}
		tUTC := t0.UTC().Truncate(24 * time.Hour)
		fromPtr = &tUTC
	}
	if f := strings.TrimSpace(c.Query("to")); f != "" {
		t1, err := time.Parse("2006-01-02", f)
		if err != nil {
			c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse(err.Error()))
			return
		}
		tUTC := t1.UTC().Truncate(24 * time.Hour)
		toPtr = &tUTC
	}

	out, err := h.summaryUC.Execute(c.Request.Context(), expensesusecases.SummaryInput{
		ProjectID: pid,
		UserID:    c.GetString(userhttp.AuthUserIDKey),
		UserRole:  c.GetString(userhttp.AuthRoleKey),
		From:      fromPtr,
		To:        toPtr,
	})
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, summaryToResp(out))
}

func (h *Handler) Get(c *gin.Context) {
	eid := c.Param("id")
	exp, err := h.getExp.Execute(c.Request.Context(), expensesusecases.GetExpenseInput{
		ID:       eid,
		UserID:   c.GetString(userhttp.AuthUserIDKey),
		UserRole: c.GetString(userhttp.AuthRoleKey),
	})
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, expenseToResp(*exp, ""))
}

func (h *Handler) Patch(c *gin.Context) {
	var body PatchExpenseBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse(err.Error()))
		return
	}

	input := expensesusecases.UpdateExpenseInput{
		ID:       c.Param("id"),
		ActorID:  c.GetString(userhttp.AuthUserIDKey),
		UserRole: c.GetString(userhttp.AuthRoleKey),
	}

	if body.Amount != nil {
		input.AmountStr = body.Amount
	}
	if body.Description != nil {
		input.Description = body.Description
	}
	if body.Currency != nil {
		input.Currency = body.Currency
	}
	if body.ReceiptStoragePath != nil {
		input.ReceiptPath = body.ReceiptStoragePath
	}
	if body.ExpenseDate != nil {
		ed, err := parseExpenseDate("2006-01-02", *body.ExpenseDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse(err.Error()))
			return
		}
		input.ExpenseDate = &ed
	}

	exp, err := h.upd.Execute(c.Request.Context(), input)
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, expenseToResp(*exp, ""))
}

func (h *Handler) Approve(c *gin.Context) {
	err := h.approve.Execute(c.Request.Context(), expensesusecases.MutationInput{
		ExpenseID: c.Param("id"),
		ActorID:   c.GetString(userhttp.AuthUserIDKey),
		UserRole:  c.GetString(userhttp.AuthRoleKey),
	})
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) Delete(c *gin.Context) {
	err := h.deleteUC.Execute(c.Request.Context(), expensesusecases.DeleteExpenseInput{
		ExpenseID: c.Param("id"),
		ActorID:   c.GetString(userhttp.AuthUserIDKey),
		UserRole:  c.GetString(userhttp.AuthRoleKey),
	})
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) Reject(c *gin.Context) {
	var body RejectBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse(err.Error()))
		return
	}

	err := h.reject.Execute(c.Request.Context(), expensesusecases.RejectExpenseInput{
		MutationInput: expensesusecases.MutationInput{
			ExpenseID: c.Param("id"),
			ActorID:   c.GetString(userhttp.AuthUserIDKey),
			UserRole:  c.GetString(userhttp.AuthRoleKey),
		},
		Reason: body.RejectionReason,
	})
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) ReceiptUpload(c *gin.Context) {
	var body ReceiptUploadBody
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

func (h *Handler) ListNotifications(c *gin.Context) {
	params := pagination.NewParams(parsePage(c.DefaultQuery("page", "1"), 1), parsePage(c.DefaultQuery("page_size", "20"), 20))
	userID := c.GetString(userhttp.AuthUserIDKey)
	result, err := h.listNotifs.Execute(c.Request.Context(), userID, params)
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, toExpenseNotificationListResponse(result))
}

func (h *Handler) MarkNotificationRead(c *gin.Context) {
	userID := c.GetString(userhttp.AuthUserIDKey)
	if err := h.markNotifRead.Execute(c.Request.Context(), c.Param("id"), userID); err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) UnreadNotificationCount(c *gin.Context) {
	userID := c.GetString(userhttp.AuthUserIDKey)
	count, err := h.unreadNotifs.Execute(c.Request.Context(), userID)
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, gin.H{"unread_count": count})
}
