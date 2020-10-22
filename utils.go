package main

import (
	"github.com/labstack/echo/v4"
)

func responseError(status int, description string, e error) error {
	Log.Println(description, e)
	return echo.NewHTTPError(status, description)
}
