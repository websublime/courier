package api

import (
	"fmt"

	"github.com/gobuffalo/uuid"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/websublime/courier/config"
	"github.com/websublime/courier/models"
	"github.com/websublime/courier/storage"
	"github.com/websublime/courier/utils"
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
	Audience   *models.Audience
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
	router.Post("audience", api.CreateAudience)
	router.Post("topic", api.CreateTopic)
	router.Get("topic", api.GetTopics)
}

func NewSocketAPI(api *API, router fiber.Router) {
	router.Use(func(ctx *fiber.Ctx) error {
		query := ctx.Query("token")
		param := utils.Decrypt([]byte(api.config.CourierKeySecret), query)

		token, err := utils.ParseJwtToken(param, api.config.CourierJWTSecret)
		if err != nil {
			return utils.NewException(utils.ErrorInvalidToken, fiber.StatusForbidden, err.Error())
		}

		claimer, ok := token.Claims.(*utils.GoTrueClaims)
		if !ok {
			return utils.NewException(utils.ErrorServerUnknown, fiber.StatusBadRequest, "Your token is not valid")
		}

		audience, err := models.FindAudience(api.db, claimer.Audience)
		if err != nil {
			return utils.NewException(utils.ErrorAudienceNotFound, fiber.StatusBadRequest, err.Error())
		}

		if websocket.IsWebSocketUpgrade(ctx) {
			ctx.Locals("token", token)
			ctx.Locals("audience", audience)
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
			audience := connection.Locals("audience").(*models.Audience)
			uid, _ := uuid.FromString(fmt.Sprintf("%s", connection.Locals("requestid")))

			api.clients = append(api.clients, Client{
				ID:         uid,
				Audience:   audience,
				Connection: connection,
			})

			payload := []byte(uid.String())

			connection.WriteMessage(1, payload)
		}
	}
}
