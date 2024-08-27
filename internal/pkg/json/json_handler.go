package json

import (
	"errors"
	"mlvt/internal/infra/zap-logging/log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// JSONResponse defines the structure of a typical JSON response
type JSONResponse struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// WriteJSON writes a JSON response with the given status and data.
func WriteJSON(ctx *gin.Context, status int, data interface{}) {
	// Construct the response body
	respBody := JSONResponse{
		Error:   false,
		Message: "Success",
		Data:    data,
	}

	// Write the JSON response
	ctx.JSON(status, respBody)

	log.Info("Response written ", "status: ", status, " data: ", data)
}

// ReadJSON reads and binds a JSON request body to a struct.
func ReadJSON(ctx *gin.Context, data interface{}) error {
	// Attempt to bind JSON data to the provided struct
	if err := ctx.ShouldBindJSON(data); err != nil {
		log.Errorf("Failed to bind JSON: %s", err.Error())
		return errors.New("request format error: " + err.Error())
	}

	log.Info("JSON request body read successfully ", "data: ", data)
	return nil
}

// ErrorJSON sends a JSON error response.
func ErrorJSON(ctx *gin.Context, err error, status ...int) {
	// Default to 400 Bad Request if no status is provided
	statusCode := http.StatusBadRequest
	if len(status) > 0 {
		statusCode = status[0]
	}

	log.Error("Error occurred ", "error: ", err.Error())

	// Create a response payload
	payload := JSONResponse{
		Error:   true,
		Message: err.Error(),
	}

	// Write the JSON response
	ctx.JSON(statusCode, payload)
}
