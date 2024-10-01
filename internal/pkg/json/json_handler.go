package json

import (
	"errors"
	"mlvt/internal/infra/reason"
	"mlvt/internal/infra/zap-logging/log"

	"github.com/gin-gonic/gin"
)

// JSONResponse defines the structure of a typical JSON response
type JSONResponse struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// WriteJSON writes a JSON response with the given status and data.
func WriteJSON(ctx *gin.Context, status int, message string, data interface{}) {
	// Construct the response body
	respBody := JSONResponse{
		Error:   status >= 400,
		Message: message,
		Data:    data,
	}

	// Write the JSON response
	ctx.JSON(status, respBody)
	log.Infof("%s--> %s : %d - %s : %v", reason.ResponseWritten.Message(), reason.Status.Message(), status, reason.Data.Message(), data)
}

// ReadJSON reads and binds a JSON request body to a struct.
func ReadJSON(ctx *gin.Context, data interface{}) error {
	if err := ctx.ShouldBindJSON(data); err != nil {
		log.Errorf(reason.FailedToBindJSON.Message()+": %s", err.Error())
		return errors.New(reason.RequestFormatError.Message() + ": " + err.Error())
	}
	return nil
}

// ErrorJSON sends a JSON error response.
func ErrorJSON(ctx *gin.Context, message string, status int) {
	log.Error("Error occurred", "error: ", message)
	ctx.JSON(status, JSONResponse{
		Error:   true,
		Message: message,
	})
}
