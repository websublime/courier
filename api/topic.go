package api

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/websublime/courier/models"
	"github.com/websublime/courier/storage"
	"github.com/websublime/courier/utils"
)

func (api *API) CreateTopic(ctx *fiber.Ctx) error {
	token := ctx.Locals("token").(*jwt.Token)
	claimer, ok := token.Claims.(*utils.GoTrueClaims)

	topic := new(models.Topic)

	if !ok {
		return utils.NewException(utils.ErrorServerUnknown, fiber.StatusBadRequest, "Your token is not valid")
	}

	audience, err := models.FindAudience(api.db, claimer.Audience)
	if err != nil {
		return utils.NewException(utils.ErrorAudienceNotFound, fiber.StatusBadRequest, err.Error())
	}

	if err := ctx.BodyParser(topic); err != nil {
		return utils.NewException(utils.ErrorTopicsInvalid, fiber.StatusPreconditionFailed, err.Error())
	}

	topic.Audience = audience
	topic.AudienceID = audience.ID

	err = api.db.Transaction(func(tx *storage.Connection) error {
		terr := tx.Create(topic)

		return terr
	})
	if err != nil {
		return utils.NewException(utils.ErrorTopicCreate, fiber.StatusBadRequest, err.Error())
	}

	return ctx.JSON(fiber.Map{
		"data": topic,
	})
}

func (api *API) GetTopics(ctx *fiber.Ctx) error {
	token := ctx.Locals("token").(*jwt.Token)
	claimer, ok := token.Claims.(*utils.GoTrueClaims)

	if !ok {
		return utils.NewException(utils.ErrorServerUnknown, fiber.StatusBadRequest, "Your token is not valid")
	}

	audience, err := models.FindAudience(api.db, claimer.Audience)
	if err != nil {
		return utils.NewException(utils.ErrorAudienceNotFound, fiber.StatusBadRequest, err.Error())
	}

	topics, err := models.FindTopicsByAudienceID(api.db, audience.ID)
	if err != nil {
		return utils.NewException(utils.ErrorTopicNotFound, fiber.StatusBadRequest, err.Error())
	}

	return ctx.JSON(fiber.Map{
		"data": topics,
	})
}
