package api

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/websublime/courier/models"
	"github.com/websublime/courier/storage"
	"github.com/websublime/courier/utils"
)

func (api *API) CreateAudience(ctx *fiber.Ctx) error {
	token := ctx.Locals("token").(*jwt.Token)
	claimer, ok := token.Claims.(*utils.GoTrueClaims)

	if !ok {
		return utils.NewException(utils.ErrorServerUnknown, fiber.StatusBadRequest, "Your token is not valid")
	}

	audience, err := models.NewAudience(claimer.Audience)
	if err != nil {
		return utils.NewException(utils.ErrorAudienceCreate, fiber.StatusBadRequest, err.Error())
	}

	err = api.db.Transaction(func(tx *storage.Connection) error {
		terr := tx.Create(audience)

		return terr
	})
	if err != nil {
		return utils.NewException(utils.ErrorAudienceCreate, fiber.StatusBadRequest, err.Error())
	}

	return ctx.JSON(fiber.Map{
		"data": audience,
	})
}
