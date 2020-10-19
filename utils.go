package main

import (
	"fmt"
	"github.com/labstack/echo"
)

func responseError(status int, description string, e error) error {
	w := fmt.Errorf("%s: %v", description, e)
	Log.Println(w)
	return echo.NewHTTPError(status, w.Error())
}
