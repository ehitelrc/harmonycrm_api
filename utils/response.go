package utils

import (
	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Error   string      `json:"error,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func Respond(c *gin.Context, status int, success bool, message string, data interface{}, err error) {
	response := APIResponse{
		Success: success,
		Message: message,
	}

	if err != nil {
		response.Error = err.Error()
	} else {
		response.Data = data
	}

	c.JSON(status, response)
}
