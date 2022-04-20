package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

const (
	APIVersion = "v1.0"
)

var (
	APIEnv   string
	APIPort  string
	APIHTTPS bool
)

const (
	DefaultEnv  = "dev"
	DefaultPort = "1323"
)

const (
	AdminUser = "admin"
	AdminPass = "passw0rd"
)

func authUser(u, p string) bool {
	return u == AdminUser && p == AdminPass
}

func init() {
	APIEnv = os.Getenv("API_ENV")
	if APIEnv == "" {
		APIEnv = DefaultEnv
	}
	APIPort = os.Getenv("API_PORT")
	if APIPort == "" {
		APIPort = DefaultPort
	}
	APIHTTPS = len(os.Getenv("API_HTTPS")) > 0
}

func main() {
	// Create the API
	e := echo.New()
	e.Logger.SetLevel(log.INFO)
	e.Logger.Info("API created")

	// Configure the middleware
	e.HideBanner = true
	e.HidePort = true
	e.Use(middleware.Logger())
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.RequestID())

	// Add the routes
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

	e.POST("/auth", func(c echo.Context) error {
		hs := c.Request().Header
		user := hs.Get("X-Auth-User")
		pass := hs.Get("X-Auth-Pass")

		if user == AdminUser && pass == AdminPass {
			return c.JSON(http.StatusOK, map[string]any{
				"success":    true,
				"apiVersion": APIVersion,
				"message":    "Authenticated",
			})
		}

		return c.JSON(http.StatusUnauthorized, map[string]any{
			"success":    false,
			"apiVersion": APIVersion,
			"message":    "Unauthorized",
		})
	})

	e.GET("/error/status", func(c echo.Context) error {
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"success":    false,
			"apiVersion": APIVersion,
			"message":    "Something went wrong!",
		})
	})
	e.GET("/error/return", func(c echo.Context) error {
		return errors.New("something went wrong: can't read something from somewhere")
	})
	e.GET("/error/panic", func(c echo.Context) error {
		panic(errors.New("something went wrong: can't read something from somewhere"))
		// return nil
	})

	// Start the API server
	go func() {
		e.Logger.Info("API starting")
		if err := e.Start(":1323"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal()
		}
	}()

	// Watch for cancelation
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	s := <-quit
	e.Logger.Info(fmt.Sprintf("Got shutdown signal: %s", s))

	// Give the API a chance to shutdown gracefully
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}

	e.Logger.Info("Done.")
}
