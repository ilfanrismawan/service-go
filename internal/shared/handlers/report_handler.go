package handlers

import (
	"net/http"
	"service/internal/shared/model"
	service "service/internal/shared/service"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// ReportHandler handles report endpoints
type ReportHandler struct {
	reportService *service.ReportService
}

// NewReportHandler creates a new report handler
func NewReportHandler() *ReportHandler {
	return &ReportHandler{
		reportService: service.NewReportService(),
	}
}

// GetMonthlyReport godoc
// @Summary Get monthly report
// @Description Get comprehensive monthly report for a specific month and year
// @Tags reports
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param year query int true "Year" default(2024)
// @Param month query int true "Month (1-12)" default(1)
// @Success 200 {object} core.APIResponse
// @Failure 400 {object} core.ErrorResponse
// @Failure 401 {object} core.ErrorResponse
// @Failure 500 {object} core.ErrorResponse
// @Router /reports/monthly [get]
func (h *ReportHandler) GetMonthlyReport(c *gin.Context) {
	// Parse query parameters
	yearStr := c.DefaultQuery("year", strconv.Itoa(time.Now().Year()))
	monthStr := c.DefaultQuery("month", strconv.Itoa(int(time.Now().Month())))

	year, err := strconv.Atoi(yearStr)
	if err != nil || year < 2020 || year > 2030 {
		c.JSON(http.StatusBadRequest, core.CreateErrorResponse(
			"validation_error",
			"Invalid year parameter",
			nil,
		))
		return
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil || month < 1 || month > 12 {
		c.JSON(http.StatusBadRequest, core.CreateErrorResponse(
			"validation_error",
			"Invalid month parameter (must be 1-12)",
			nil,
		))
		return
	}

	report, err := h.reportService.GenerateMonthlyReport(c.Request.Context(), year, month)
	if err != nil {
		c.JSON(http.StatusInternalServerError, core.CreateErrorResponse(
			"report_generation_failed",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, core.SuccessResponse(report, "Monthly report generated successfully"))
}

// GetYearlyReport godoc
// @Summary Get yearly report
// @Description Get comprehensive yearly report for a specific year
// @Tags reports
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param year query int true "Year" default(2024)
// @Success 200 {object} core.APIResponse
// @Failure 400 {object} core.ErrorResponse
// @Failure 401 {object} core.ErrorResponse
// @Failure 500 {object} core.ErrorResponse
// @Router /reports/yearly [get]
func (h *ReportHandler) GetYearlyReport(c *gin.Context) {
	// Parse query parameters
	yearStr := c.DefaultQuery("year", strconv.Itoa(time.Now().Year()))

	year, err := strconv.Atoi(yearStr)
	if err != nil || year < 2020 || year > 2030 {
		c.JSON(http.StatusBadRequest, core.CreateErrorResponse(
			"validation_error",
			"Invalid year parameter",
			nil,
		))
		return
	}

	report, err := h.reportService.GetYearlyReport(c.Request.Context(), year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, core.CreateErrorResponse(
			"report_generation_failed",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, core.SuccessResponse(report, "Yearly report generated successfully"))
}

// GetCurrentMonthReport godoc
// @Summary Get current month report
// @Description Get report for the current month
// @Tags reports
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} core.APIResponse
// @Failure 401 {object} core.ErrorResponse
// @Failure 500 {object} core.ErrorResponse
// @Router /reports/current-month [get]
func (h *ReportHandler) GetCurrentMonthReport(c *gin.Context) {
	now := time.Now()
	report, err := h.reportService.GenerateMonthlyReport(c.Request.Context(), now.Year(), int(now.Month()))
	if err != nil {
		c.JSON(http.StatusInternalServerError, core.CreateErrorResponse(
			"report_generation_failed",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, core.SuccessResponse(report, "Current month report generated successfully"))
}

// GetReportSummary godoc
// @Summary Get report summary
// @Description Get summary of reports for the last 12 months
// @Tags reports
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} core.APIResponse
// @Failure 401 {object} core.ErrorResponse
// @Failure 500 {object} core.ErrorResponse
// @Router /reports/summary [get]
func (h *ReportHandler) GetReportSummary(c *gin.Context) {
	now := time.Now()
	var summary []core.MonthlyReportData

	// Get data for the last 12 months
	for i := 11; i >= 0; i-- {
		date := now.AddDate(0, -i, 0)
		report, err := h.reportService.GenerateMonthlyReport(c.Request.Context(), date.Year(), int(date.Month()))
		if err != nil {
			continue
		}

		summary = append(summary, core.MonthlyReportData{
			Month:   report.Month,
			Orders:  report.TotalOrders,
			Revenue: report.TotalRevenue,
		})
	}

	c.JSON(http.StatusOK, core.SuccessResponse(summary, "Report summary generated successfully"))
}
