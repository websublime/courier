package api

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/websublime/courier/models"
	"github.com/websublime/courier/storage"
	"github.com/websublime/courier/utils"
)

func (api *API) CreateNewChannel(ctx *fiber.Ctx) error {
	token := ctx.Locals("token").(*jwt.Token)
	claimer := token.Claims.(*utils.TrueClaims)

	aud, err := models.FindAudience(api.db, claimer.Audience)
	if err != nil {
		return utils.NewException(utils.ErrorAudienceNotFound, fiber.StatusNotFound, err.Error())
	}

	channel, err := models.NewChannel("")

	if err := ctx.BodyParser(channel); err != nil {
		return utils.NewException(utils.ErrorBodyParse, fiber.StatusPreconditionFailed, err.Error())
	}

	channel.Audience = aud
	channel.AudienceID = aud.ID

	err = api.db.Transaction(func(tx *storage.Connection) error {
		terr := tx.Create(channel)

		return terr
	})
	if err != nil {
		return utils.NewException(utils.ErrorChannelCreation, fiber.StatusBadRequest, err.Error())
	}

	return ctx.JSON(fiber.Map{
		"data": channel,
	})
}
