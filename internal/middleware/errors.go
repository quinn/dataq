package middleware

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.quinn.io/dataq/ui"
)

type ToastData struct {
	Message string `json:"message"`
	Type    string `json:"type"`
}

func HTTPErrorHandler(err error, c echo.Context) {
	statusCode := http.StatusInternalServerError
	if httpErr, ok := err.(*echo.HTTPError); ok {
		statusCode = httpErr.Code
	}

	var renderErr error
	if c.Request().Header.Get("HX-Request") == "true" {
		c.Response().Header().Set("HX-Retarget", "body")
		c.Response().Header().Set("HX-Reswap", "beforeend")

		eventData := map[string]interface{}{
			"showToast": ToastData{
				Message: err.Error(),
				Type:    "error",
			},
		}

		event, _ := json.Marshal(eventData)

		c.Response().Header().Set("HX-Trigger", string(event))
	} else {
		c.Response().WriteHeader(statusCode)
		renderErr = ui.ErrorPage(err.Error()).Render(c.Request().Context(), c.Response().Writer)
	}

	if renderErr != nil {
		_ = c.JSON(statusCode, err)
	}
}
