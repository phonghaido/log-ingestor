package helpers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type APIError struct {
	StatusCode int    `json:"statusCode"`
	Msg        string `json:"msg"`
}

func (e APIError) Error() string {
	return fmt.Sprintf("api error message: %d", e.StatusCode)
}

func NewAPIError(statusCode int, err error) APIError {
	return APIError{
		StatusCode: statusCode,
		Msg:        err.Error(),
	}
}

func InvalidJSON(c echo.Context) error {
	return NewAPIError(http.StatusBadRequest, fmt.Errorf("invalid json data request"))
}

type APIFunc func(c echo.Context) error

func ErrorWrapper(h APIFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := h(c); err != nil {
			if apiErr, ok := err.(APIError); ok {
				return WriteJSON(c, apiErr.StatusCode, apiErr)
			} else {
				errResp := map[string]any{
					"statusCode": http.StatusInternalServerError,
					"msg":        "internal server error",
				}
				return WriteJSON(c, http.StatusInternalServerError, errResp)
			}
		}
		return nil
	}
}

func WriteJSON(c echo.Context, statusCode int, v any) error {
	return c.JSON(statusCode, v)
}
