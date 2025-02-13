package order

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"net/http"
)

var (
	validate *validator.Validate // use a singleton to be thread safe
)

func init() {
	validate = validator.New()
}

// Validate runs through the tags and validates all fields.
func Validate(model interface{}) error {
	if err := validate.Struct(model); err != nil {
		e, ok := err.(validator.ValidationErrors)
		if ok {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("error validating model: %s", e.Error()))
		}
		return err
	}
	return nil
}
