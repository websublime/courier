package config

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/websublime/courier/utils"
)

// BootApplication initializes fiber application
func BootApplication(conf *EnvironmentConfig) *fiber.App {
	app := fiber.New(fiber.Config{
		ServerHeader: "courier-storage",
		Prefork:      conf.CourierProduction,
		ErrorHandler: ErrorHandler,
	})

	app.Use(recover.New())
	app.Use(pprof.New())
	app.Use(compress.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
		AllowMethods: "GET,POST",
	}))
	app.Use(etag.New())
	app.Use(logger.New())
	app.Use(favicon.New())
	app.Use(requestid.New())

	return app
}

// ErrorHandler handle fiber errors
func ErrorHandler(ctx *fiber.Ctx, err error) error {
	// Statuscode defaults to 500
	code := fiber.StatusInternalServerError
	message := "Internal server error"
	exception := utils.ErrorServerUnknown

	// Retreive the custom statuscode if it's an fiber.*Error
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}

	if e, ok := err.(*utils.Exception); ok {
		code = e.Code
		message = e.Message
		exception = e.Exception
	}

	// Send custom error page
	ctx.Set(fiber.HeaderContentType, "application/json")

	return ctx.Status(code).JSON(fiber.Map{
		"status": code,
		"error":  message,
		"code":   exception,
	})
}
