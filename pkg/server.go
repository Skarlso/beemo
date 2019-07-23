package pkg

import (
	"log"

	"github.com/Skarlso/acquia-beemo/internal"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Serve will start the echo server listening for webhooks.
func Serve() error {
	log.Println("Starting listener...")
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	labeler := internal.NewGithubLabeler()
	// Routes
	e.POST("/githook", GitWebHook(labeler))

	// Start server
	return e.Start(":9998")
}
