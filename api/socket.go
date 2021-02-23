package api

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gofiber/websocket/v2"
	"github.com/gofrs/uuid"
	"github.com/websublime/courier/models"
	"github.com/websublime/courier/storage"
	"github.com/websublime/courier/utils"
)

const (
	PUBLISH     = "publish"
	BROADCAST   = "broadcast"
	SUBSCRIBE   = "subscribe"
	UNSUBSCRIBE = "unsubscribe"
)

func (api *API) SocketHandler(ctx *websocket.Conn) {
	api.HandleReceivedMessage(ctx)
}

func (api *API) HandleReceivedMessage(ctx *websocket.Conn) {
	defer func() {
		api.unregisterClient <- ctx
		ctx.Close()
	}()

	api.registerClient <- ctx

	for {
		message := utils.Message{}

		messageType, payload, err := ctx.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println("read error:", err)
			}

			return
		}

		err = json.Unmarshal(payload, &message)
		if err != nil {
			log.Println("Payload is invalid")
			return
		}

		switch message.Action {
		case PUBLISH:
			api.Publish(ctx, message.Topic, message.Message)
			api.SaveMessage(ctx, message.Topic, message.Message)
			log.Println("publisged to", message.Topic)
			break
		case SUBSCRIBE:
			api.Subscribe(ctx, message.Topic)
			log.Println("new subscription to topic", message.Topic, len(api.subscriptions), "subscribed")
			break
		case UNSUBSCRIBE:
			api.Unsubscribe(ctx, message.Topic)
			log.Println("Client want to unsubscribe the topic", message.Topic)
			break
		case BROADCAST:
			api.Broadcast(ctx, message.Message)
			log.Println("broadcast to", len(api.subscriptions), "subscribed")
			break
		}

		log.Println("websocket message received of type", messageType)
		log.Println("websocket message", message)
	}
}

func (api *API) Subscribe(ctx *websocket.Conn, topic string) {
	uid, _ := uuid.FromString(fmt.Sprintf("%s", ctx.Locals("requestid")))
	audience := ctx.Locals("audience").(*models.Audience)

	var client Client
	var subscriptionList []Subscription

	for _, cl := range api.clients {
		if cl.ID == uid && cl.Audience.ID == audience.ID {
			client = cl
		}
	}

	for _, subscription := range api.subscriptions {
		if subscription.Client.ID == client.ID && subscription.Topic == topic {
			subscriptionList = append(subscriptionList, subscription)
		}
	}

	if len(subscriptionList) == 0 {
		api.subscriptions = append(api.subscriptions, Subscription{
			Topic:  topic,
			Client: &client,
		})
	}
}

func (api *API) Unsubscribe(ctx *websocket.Conn, topic string) {
	uid, _ := uuid.FromString(fmt.Sprintf("%s", ctx.Locals("requestid")))
	audience := ctx.Locals("audience").(*models.Audience)

	var client Client

	for _, cl := range api.clients {
		if cl.ID == uid && cl.Audience.ID == audience.ID {
			client = cl
		}
	}

	for index, subscription := range api.subscriptions {
		if subscription.Client.ID == client.ID && subscription.Topic == topic {
			api.subscriptions = append(api.subscriptions[:index], api.subscriptions[index+1:]...)
		}
	}
}

func (api *API) Publish(ctx *websocket.Conn, topic string, message json.RawMessage) {
	audience := ctx.Locals("audience").(*models.Audience)
	var subscriptionList []Subscription

	for _, subscription := range api.subscriptions {
		if subscription.Topic == topic && subscription.Client.Audience.ID == audience.ID {
			subscriptionList = append(subscriptionList, subscription)
		}
	}

	for _, sub := range subscriptionList {
		msg, _ := message.MarshalJSON()

		sub.Client.Connection.WriteMessage(1, msg)
	}
}

func (api *API) Broadcast(ctx *websocket.Conn, message json.RawMessage) {
	audience := ctx.Locals("audience").(*models.Audience)

	for _, subscription := range api.subscriptions {
		if subscription.Client.Audience.ID == audience.ID {
			msg, _ := message.MarshalJSON()

			subscription.Client.Connection.WriteMessage(1, msg)
		}
	}
}

func (api *API) SaveMessage(ctx *websocket.Conn, topicName string, message json.RawMessage) error {
	audience := ctx.Locals("audience").(*models.Audience)

	topic, err := models.FindTopicByNameAndAudienceID(api.db, topicName, audience.ID)
	if err != nil {
		return err
	}

	msg, err := models.NewMessage(message, topic.ID)
	if err != nil {
		return err
	}

	err = api.db.Transaction(func(tx *storage.Connection) error {
		terr := tx.Create(msg)

		return terr
	})

	return err
}
