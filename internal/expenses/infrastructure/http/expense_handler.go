package expenseshttp

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	expensesports "macabi-back/internal/expenses/application/ports"
	expensesusecases "macabi-back/internal/expenses/application/usecases"
	expensesdomain "macabi-back/internal/expenses/domain"
	expensesdto "macabi-back/internal/expenses/infrastructure/http/dto"
	sharederrors "macabi-back/internal/shared/errors"
	"macabi-back/internal/shared/pagination"
	userhttp "macabi-back/internal/user/infrastructure/http"
)

func (h *Handler) Create(c *gin.Context) {
	if isMultipartRequest(c) {
		h.createMultipart(c)
		return
	}
	var body expensesdto.CreateExpenseBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse(err.Error()))
		return
	}
	exp, err := h.createExpenseJSON(c, body)
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusCreated, expensesdto.ExpenseToResp(*exp, ""))
}

func (h *Handler) createExpenseJSON(c *gin.Context, body expensesdto.CreateExpenseBody) (*expensesdomain.Expense, error) {
	d, err := expensesdto.ParseExpenseDate("2006-01-02", body.ExpenseDate)
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
		CategoryID:  body.CategoryID,
	})
}

func (h *Handler) createMultipart(c *gin.Context) {
	const maxBytes = expensesusecases.MaxReceiptBytes
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBytes+64*1024)
	if err := c.Request.ParseMultipartForm(maxBytes + 64*1024); err != nil {
		c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse("no se pudo leer el formulario"))
		return
	}
	d, err := expensesdto.ParseExpenseDate("2006-01-02", strings.TrimSpace(c.PostForm("expense_date")))
	if err != nil {
		c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse(err.Error()))
		return
	}
	base := expensesusecases.CreateExpenseInput{
		ProjectID:   strings.TrimSpace(c.PostForm("project_id")),
		SubmittedBy: c.GetString(userhttp.AuthUserIDKey),
		UserRole:    c.GetString(userhttp.AuthRoleKey),
		AmountStr:   strings.TrimSpace(c.PostForm("amount")),
		Description: strings.TrimSpace(c.PostForm("description")),
		ExpenseDate: d,
		Currency:    strings.TrimSpace(c.PostForm("currency")),
	}
	file, err := c.FormFile("file")
	if err != nil {
		exp, err := h.create.Execute(c.Request.Context(), base)
		if err != nil {
			c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
			return
		}
		c.JSON(http.StatusCreated, expensesdto.ExpenseToResp(*exp, ""))
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
	c.JSON(http.StatusCreated, expensesdto.ExpenseToResp(*exp, ""))
}

func (h *Handler) Get(c *gin.Context) {
	exp, err := h.getExp.Execute(c.Request.Context(), expensesusecases.GetExpenseInput{
		ID:       c.Param("id"),
		UserID:   c.GetString(userhttp.AuthUserIDKey),
		UserRole: c.GetString(userhttp.AuthRoleKey),
	})
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, expensesdto.DetailItemToExpenseResponse(*exp))
}

func (h *Handler) Patch(c *gin.Context) {
	var body expensesdto.PatchExpenseBody
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
	if body.CategoryID != nil {
		input.CategoryID = body.CategoryID
	}
	if body.ExpenseDate != nil {
		ed, err := expensesdto.ParseExpenseDate("2006-01-02", *body.ExpenseDate)
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
	c.JSON(http.StatusOK, expensesdto.ExpenseToResp(*exp, ""))
}

func (h *Handler) Approve(c *gin.Context) {
	if err := h.approve.Execute(c.Request.Context(), expensesusecases.MutationInput{
		ExpenseID: c.Param("id"),
		ActorID:   c.GetString(userhttp.AuthUserIDKey),
		UserRole:  c.GetString(userhttp.AuthRoleKey),
	}); err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) Reject(c *gin.Context) {
	var body expensesdto.RejectBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse(err.Error()))
		return
	}
	if err := h.reject.Execute(c.Request.Context(), expensesusecases.RejectExpenseInput{
		MutationInput: expensesusecases.MutationInput{
			ExpenseID: c.Param("id"),
			ActorID:   c.GetString(userhttp.AuthUserIDKey),
			UserRole:  c.GetString(userhttp.AuthRoleKey),
		},
		Reason: body.RejectionReason,
	}); err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) Delete(c *gin.Context) {
	if err := h.deleteUC.Execute(c.Request.Context(), expensesusecases.DeleteExpenseInput{
		ExpenseID: c.Param("id"),
		ActorID:   c.GetString(userhttp.AuthUserIDKey),
		UserRole:  c.GetString(userhttp.AuthRoleKey),
	}); err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) ListMine(c *gin.Context) {
	params := pagination.ParseParams(c.Query("page"), c.Query("page_size"))
	res, err := h.listMine.Execute(c.Request.Context(), c.GetString(userhttp.AuthUserIDKey), params)
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, expensesdto.ListResp(res))
}

func (h *Handler) ListAll(c *gin.Context) {
	params := pagination.ParseParams(c.Query("page"), c.Query("page_size"))
	filter := expensesports.ExpenseListFilter{
		ProjectID: strings.TrimSpace(c.Query("project_id")),
		Status:    strings.TrimSpace(c.Query("status")),
		Query:     strings.TrimSpace(c.Query("q")),
	}
	if f := strings.TrimSpace(c.Query("from")); f != "" {
		t0, err := time.Parse("2006-01-02", f)
		if err != nil {
			c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse(err.Error()))
			return
		}
		tUTC := t0.UTC().Truncate(24 * time.Hour)
		filter.From = &tUTC
	}
	if f := strings.TrimSpace(c.Query("to")); f != "" {
		t1, err := time.Parse("2006-01-02", f)
		if err != nil {
			c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse(err.Error()))
			return
		}
		tUTC := t1.UTC().Truncate(24 * time.Hour)
		filter.To = &tUTC
	}
	res, err := h.listAll.Execute(c.Request.Context(), filter, params)
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, expensesdto.ListResp(res))
}

func (h *Handler) ListByProject(c *gin.Context) {
	params := pagination.ParseParams(c.Query("page"), c.Query("page_size"))
	res, err := h.listProj.Execute(c.Request.Context(), expensesusecases.ListProjectExpensesInput{
		ProjectID: c.Param("id"),
		UserID:    c.GetString(userhttp.AuthUserIDKey),
		UserRole:  c.GetString(userhttp.AuthRoleKey),
		Params:    params,
	})
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, expensesdto.ListResp(res))
}
