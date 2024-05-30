package helpers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/phonghaido/log-ingestor/types"
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

func InvalidJSON() error {
	return NewAPIError(http.StatusBadRequest, fmt.Errorf("invalid json data request"))
}

func ValidateLogData(logData types.LogData) error {
	if logData.Level == "" {
		return NewAPIError(http.StatusBadRequest, fmt.Errorf("missing data 'level'"))
	}

	if logData.Commit == "" {
		return NewAPIError(http.StatusBadRequest, fmt.Errorf("missing data 'commit'"))
	}

	if logData.Message == "" {
		return NewAPIError(http.StatusBadRequest, fmt.Errorf("missing data 'message'"))
	}

	if logData.ResourceID == "" {
		return NewAPIError(http.StatusBadRequest, fmt.Errorf("missing data 'resource_id'"))
	}

	if logData.SpanID == "" {
		return NewAPIError(http.StatusBadRequest, fmt.Errorf("missing data 'span_id'"))
	}

	if logData.Timestamp == "" {
		return NewAPIError(http.StatusBadRequest, fmt.Errorf("missing data 'timestamp'"))
	} else if _, err := time.Parse(time.RFC3339, logData.Timestamp); err != nil {
		return NewAPIError(http.StatusBadRequest, fmt.Errorf("timestamp should be in RFC3339 format"))
	}

	if logData.TraceID == "" {
		return NewAPIError(http.StatusBadRequest, fmt.Errorf("missing data 'trace_id'"))
	}

	if logData.Metadata.ParentResourceID == "" {
		return NewAPIError(http.StatusBadRequest, fmt.Errorf("missing data 'metadata parent_resource_id'"))
	}

	return nil
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
