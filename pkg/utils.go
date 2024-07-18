package utils

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type JSONRespone struct {
	Error   bool        `json:"error"`
	Message string      `json:"message`
	Data    interface{} `json:"data, omitempty`
}

func WriteJSON(c *gin.Context, status int, data interface{}, headers ...http.Header) error {
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for key, value := range headers[0] {
			c.Writer.Header()[key] = value
		}
	}

	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(status)

	_, err = c.Writer.Write(out)
	if err != nil {
		return err
	}

	return nil
}

func ReadJSON(c *gin.Context, data interface{}) error {
	maxBytes := 1024 * 1024 // 1 megabyte
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, int64(maxBytes))

	dec := json.NewDecoder(c.Request.Body)

	dec.DisallowUnknownFields()

	err := dec.Decode(data)
	if err != nil {
		return err
	}

	//Make sure system only receives one JSON per request
	err = dec.Decode(struct{}{})
	if err != io.EOF {
		return errors.New("body must only contain a single JSON value.")
	}

	return nil
}

func ErrorJSON(c *gin.Context, err error, status ...int) error {
	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}

	var payload JSONRespone
	payload.Error = true
	payload.Message = err.Error()

	return WriteJSON(c, statusCode, payload)
}
