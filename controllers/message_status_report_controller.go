package controllers

import (
	"harmony_api/config"
	"harmony_api/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type MessageStatusReportController struct{}

func NewMessageStatusReportController() *MessageStatusReportController {
	return &MessageStatusReportController{}
}

type StatusCount struct {
	Status string `json:"status"`
	Count  int64  `json:"count"`
}

type FailedMessageDetail struct {
	ID               uint64    `json:"id"`
	CaseID           uint      `json:"case_id"`
	ChannelMessageID string    `json:"channel_message_id"`
	TextContent      string    `json:"text_content"`
	MessageError     string    `json:"message_error"`
	CreatedAt        time.Time `json:"created_at"`
}

func (c *MessageStatusReportController) GetSummary(ctx *gin.Context) {
	companyIDParam := ctx.Param("company_id")
	companyID, err := strconv.ParseInt(companyIDParam, 10, 64)
	if err != nil {
		utils.Respond(ctx, http.StatusBadRequest, false, "company_id inválido", nil, err)
		return
	}

	// 1. Get counts grouped by status in the last 24 hours for the company
	var counts []StatusCount
	countQuery := `
		SELECT m.status, COUNT(m.id) as count
		FROM messages m
		JOIN cases c ON m.case_id = c.id
		WHERE c.company_id = ?
		  AND m.sender_type = 'agent'
		  AND m.created_at >= NOW() - INTERVAL '24 hours'
		GROUP BY m.status
	`
	if err := config.DB.Raw(countQuery, companyID).Scan(&counts).Error; err != nil {
		utils.Respond(ctx, http.StatusInternalServerError, false, "Error al obtener estadísticas de mensajes", nil, err)
		return
	}

	// 2. Get details of failed messages in the last 24 hours for the company
	var failedDetails []FailedMessageDetail
	detailsQuery := `
		SELECT m.id, m.case_id, m.channel_message_id, m.text_content, m.message_error, m.created_at
		FROM messages m
		JOIN cases c ON m.case_id = c.id
		WHERE c.company_id = ?
		  AND m.sender_type = 'agent'
		  AND (m.status = 'failed' OR m.has_error = true)
		  AND m.created_at >= NOW() - INTERVAL '24 hours'
		ORDER BY m.id DESC
	`
	if err := config.DB.Raw(detailsQuery, companyID).Scan(&failedDetails).Error; err != nil {
		utils.Respond(ctx, http.StatusInternalServerError, false, "Error al obtener detalles de fallidos", nil, err)
		return
	}

	// Map to response structure
	readCount := int64(0)
	deliveredCount := int64(0)
	sentCount := int64(0)
	failedCount := int64(0)

	for _, ct := range counts {
		switch ct.Status {
		case "read":
			readCount = ct.Count
		case "delivered":
			deliveredCount = ct.Count
		case "sent":
			sentCount = ct.Count
		case "failed":
			failedCount = ct.Count
		}
	}

	// Fallback to list length if failedCount is zero but we have failedDetails
	if failedCount == 0 && len(failedDetails) > 0 {
		failedCount = int64(len(failedDetails))
	}

	response := map[string]interface{}{
		"total":          readCount + deliveredCount + sentCount + failedCount,
		"read":           readCount,
		"delivered":      deliveredCount,
		"sent":           sentCount,
		"failed":         failedCount,
		"failed_details": failedDetails,
	}

	utils.Respond(ctx, http.StatusOK, true, "Resumen de estados de mensajes obtenido correctamente", response, nil)
}
