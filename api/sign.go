package api

import (
	"net/url"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/websublime/courier/models"
	"github.com/websublime/courier/utils"
)

func (api *API) GetSignedUrl(ctx *fiber.Ctx) error {
	token := ctx.Locals("token").(*jwt.Token)
	claimer, ok := token.Claims.(*utils.GoTrueClaims)

	if !ok {
		return utils.NewException(utils.ErrorServerUnknown, fiber.StatusBadRequest, "Your token is not valid")
	}

	_, err := models.FindAudience(api.db, claimer.Audience)
	if err != nil {
		return utils.NewException(utils.ErrorAudienceNotFound, fiber.StatusBadRequest, err.Error())
	}

	u, err := url.Parse(api.config.CourierURL)
	if err != nil {
		return utils.NewException(utils.ErrorEncryptFailure, fiber.StatusNotFound, err.Error())
	}

	signed := utils.Encrypt([]byte(api.config.CourierKeySecret), token.Raw)
	query := u.Query()
	query.Add("token", signed)
	u.RawQuery = query.Encode()

	return ctx.JSON(fiber.Map{
		"data": fiber.Map{
			"key": signed,
			"url": u.String(),
		},
	})
}
