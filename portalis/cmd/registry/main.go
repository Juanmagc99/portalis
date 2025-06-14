package main

import (
	"github.com/Juanmagc99/portalis/pkg/httperror"
	"github.com/Juanmagc99/portalis/pkg/validation"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.HTTPErrorHandler = httperror.HandleError
	e.Validator = &validation.CustomValidator{Validator: validator.New()}
}
