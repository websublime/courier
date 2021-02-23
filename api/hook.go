package api

import (
	"fmt"
	"net/url"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/websublime/courier/models"
	"github.com/websublime/courier/utils"
	"golang.org/x/net/websocket"
)

func (api *API) Hook(ctx *fiber.Ctx) error {
	token := ctx.Locals("token").(*jwt.Token)
	claimer, ok := token.Claims.(*utils.GoTrueClaims)

	hook := &models.Hook{}

	if !ok {
		return utils.NewException(utils.ErrorServerUnknown, fiber.StatusBadRequest, "Your token is not valid")
	}

	_, err := models.FindAudience(api.db, claimer.Audience)
	if err != nil {
		return utils.NewException(utils.ErrorAudienceNotFound, fiber.StatusBadRequest, err.Error())
	}

	if err := ctx.BodyParser(hook); err != nil {
		return utils.NewException(utils.ErrorHookInvalid, fiber.StatusPreconditionFailed, err.Error())
	}

	if err := hook.Validate(); len(err.Errors) > 0 {
		return utils.NewException(utils.ErrorHookInvalid, fiber.StatusUnprocessableEntity, err.Error())
	}

	u, err := url.Parse(api.config.CourierURL)
	if err != nil {
		return utils.NewException(utils.ErrorEncryptFailure, fiber.StatusNotFound, err.Error())
	}

	signed := utils.Encrypt([]byte(api.config.CourierKeySecret), token.Raw)
	query := u.Query()
	query.Add("token", signed)
	u.RawQuery = query.Encode()

	ws, err := websocket.Dial(u.String(), "", fmt.Sprintf("%s%s", ctx.BaseURL(), ctx.OriginalURL()))
	if err != nil {
		return utils.NewException(utils.ErrorSocketConnectionFailure, fiber.StatusNotFound, err.Error())
	}

	defer ws.Close()

	err = websocket.JSON.Send(ws, hook)
	if err != nil {
		return utils.NewException(utils.ErrorSocketMessageFailure, fiber.StatusNotFound, err.Error())
	}

	return ctx.JSON(fiber.Map{
		"data": hook,
	})
}
