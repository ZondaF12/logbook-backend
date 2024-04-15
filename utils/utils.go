package utils

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

var Validate = validator.New()

func ParseJSON(c echo.Context, payload any) error {
	if c.Request().Body == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid payload")
	}

	return json.NewDecoder(c.Request().Body).Decode(&payload)
}
