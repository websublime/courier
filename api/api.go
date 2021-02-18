package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/websublime/courier/config"
	"github.com/websublime/courier/storage"
	"github.com/websublime/courier/utils"
)

type API struct {
	db     *storage.Connection
	config *config.EnvironmentConfig
	app    *fiber.App
}

func WithVersion(app *fiber.App, conf *config.EnvironmentConfig, db *storage.Connection) {
	api := &API{
		db:     db,
		config: conf,
		app:    app,
	}

	router := app.Group("/v1", api.AuthorizedMiddleware)
	socketRouter := app.Group("/ws")

	NewAPI(api, router)
	NewSocketAPI(api, socketRouter)
}

func NewAPI(api *API, router fiber.Router) {
	router.Get("sign", api.GetSignedUrl)
	router.Post("channel", api.CreateNewChannel)
}

func NewSocketAPI(api *API, router fiber.Router) {
	router.Use(func(ctx *fiber.Ctx) error {
		query := ctx.Query("token")
		param := utils.Decrypt([]byte(api.config.CourierKeySecret), query)

		token, err := utils.ParseJwtToken(param, api.config.CourierJWTSecret)
		if err != nil {
			return utils.NewException(utils.ErrorInvalidToken, fiber.StatusForbidden, err.Error())
		}

		if websocket.IsWebSocketUpgrade(ctx) {
			ctx.Locals("token", token)
			return ctx.Next()
		}

		return fiber.ErrUpgradeRequired
	})
	router.Get("", websocket.New(api.SocketHandler))
}
