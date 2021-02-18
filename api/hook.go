package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/net/websocket"
)

func (api *API) Hook(ctx *fiber.Ctx) error {

	ws, err := websocket.Dial(api.config.CourierURL, "", fmt.Sprintf("%s%s", ctx.BaseURL(), ctx.OriginalURL()))
	if err != nil {
		return err
	}

	defer ws.Close()

	err = websocket.JSON.Send(ws, fiber.Map{
		"action":  "publish",
		"topic":   "system/events",
		"message": fiber.Map{"hello": "hook"},
	})

	if err != nil {
		return err
	}

	return ctx.JSON(fiber.Map{
		"data": "channel",
	})
}
