package attendancehttp

import (
	"net/http"

	"github.com/gin-gonic/gin"

	attendanceusecases "macabi-back/internal/attendance/application/usecases"
	attendancedomain "macabi-back/internal/attendance/domain"
	projectports "macabi-back/internal/project/application/ports"
	sharedErrors "macabi-back/internal/shared/errors"
	userhttp "macabi-back/internal/user/infrastructure/http"
)

type AttendanceHandler struct {
	confirmAttendance   *attendanceusecases.ConfirmAttendance
	getAttendanceCount  *attendanceusecases.GetAttendanceCount
	projectRepo         projectports.ProjectRepository
}

func NewAttendanceHandler(
	confirmAttendance *attendanceusecases.ConfirmAttendance,
	getAttendanceCount *attendanceusecases.GetAttendanceCount,
	projectRepo projectports.ProjectRepository,
) *AttendanceHandler {
	return &AttendanceHandler{
		confirmAttendance:  confirmAttendance,
		getAttendanceCount: getAttendanceCount,
		projectRepo:        projectRepo,
	}
}

func (h *AttendanceHandler) Confirm(c *gin.Context) {
	userID := c.GetString(userhttp.AuthUserIDKey)
	projectID := c.Param("id")

	input := attendanceusecases.ConfirmAttendanceInput{
		UserID:    userID,
		ProjectID: projectID,
	}

	a, err := h.confirmAttendance.Execute(c.Request.Context(), input)
	if err != nil {
		if err == attendancedomain.ErrAlreadyConfirmed {
			c.JSON(http.StatusConflict, sharedErrors.NewErrorResponse(err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, sharedErrors.NewErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, AttendanceResponse{
		ID:        a.ID,
		UserID:    a.UserID,
		ProjectID: a.ProjectID,
	})
}

func (h *AttendanceHandler) GetCount(c *gin.Context) {
	projectID := c.Param("id")

	output, err := h.getAttendanceCount.Execute(c.Request.Context(), attendanceusecases.GetAttendanceCountInput{
		ProjectID: projectID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharedErrors.NewErrorResponse(err.Error()))
		return
	}

	project, err := h.projectRepo.FindByID(c.Request.Context(), projectID)
	if err != nil {
		c.JSON(http.StatusNotFound, sharedErrors.NewErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, AttendanceCountResponse{
		ProjectID: projectID,
		Confirmed: output.Confirmed,
		Capacity:  project.Capacity,
	})
}