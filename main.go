package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

const (
	APIVersion = "v1.0"
)

var (
	APIEnv string
)

func init() {
	APIEnv = os.Getenv("API_ENV")
}

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]any{
			"success":    true,
			"apiVersion": APIVersion,
			"message":    "Hello, World!",
		})
	})
	e.GET("/greeting/:name", func(c echo.Context) error {
		name := c.Param("name")
		if name == "" {
			name = "World"
		}
		return c.JSON(http.StatusOK, map[string]any{
			"success":    true,
			"apiVersion": APIVersion,
			"message":    fmt.Sprintf("Hello, %s!", name),
		})
	})
	e.Logger.Fatal(e.Start(":1323"))
}
