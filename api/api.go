package api

import (
	"fmt"

	"github.com/gobuffalo/uuid"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/websublime/courier/config"
	"github.com/websublime/courier/storage"
)

type API struct {
	db               *storage.Connection
	config           *config.EnvironmentConfig
	app              *fiber.App
	clients          []Client
	subscriptions    []Subscription
	registerClient   chan *websocket.Conn
	unregisterClient chan *websocket.Conn
}

type Client struct {
	ID         uuid.UUID
	Connection *websocket.Conn
}

type Subscription struct {
	Topic  string
	Client *Client
}

func WithVersion(app *fiber.App, conf *config.EnvironmentConfig, db *storage.Connection) {
	api := &API{
		db:               db,
		config:           conf,
		app:              app,
		clients:          []Client{},
		subscriptions:    []Subscription{},
		registerClient:   make(chan *websocket.Conn),
		unregisterClient: make(chan *websocket.Conn),
	}

	go api.Hub()

	router := app.Group("/v1", api.AuthorizedMiddleware)
	socketRouter := app.Group("/ws")

	NewAPI(api, router)
	NewSocketAPI(api, socketRouter)
}

func NewAPI(api *API, router fiber.Router) {
	router.Get("sign", api.GetSignedUrl)
	router.Post("hook", api.Hook)
}

func NewSocketAPI(api *API, router fiber.Router) {
	router.Use(func(ctx *fiber.Ctx) error {
		/*
			query := ctx.Query("token")
			param := utils.Decrypt([]byte(api.config.CourierKeySecret), query)

			token, err := utils.ParseJwtToken(param, api.config.CourierJWTSecret)
			if err != nil {
				return utils.NewException(utils.ErrorInvalidToken, fiber.StatusForbidden, err.Error())
			}
		*/
		if websocket.IsWebSocketUpgrade(ctx) {
			// ctx.Locals("token", token)
			return ctx.Next()
		}

		return fiber.ErrUpgradeRequired
	})
	router.Get("", websocket.New(api.SocketHandler))
}

func (api *API) Hub() {
	for {
		select {
		case connection := <-api.registerClient:
			uid, _ := uuid.FromString(fmt.Sprintf("%s", connection.Locals("requestid")))
			api.clients = append(api.clients, Client{
				ID:         uid,
				Connection: connection,
			})

			payload := []byte(uid.String())

			connection.WriteMessage(1, payload)
		}
	}
}
