package api

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/websublime/courier/utils"
)

func (api *API) AuthorizedMiddleware(ctx *fiber.Ctx) error {
	auth := ctx.Get("Authorization")
	authLength := len(auth)
	authBearer := len("Bearer")

	if authLength == 0 {
		return utils.NewException(utils.ErrorStatusForbidden, fiber.StatusForbidden, "Only authorized requests are permitted")
	}

	if authLength > authBearer+1 && auth[:authBearer] == "Bearer" {
		bearer := auth[authBearer+1:]

		parser := jwt.Parser{ValidMethods: []string{jwt.SigningMethodHS256.Name}}

		token, err := parser.ParseWithClaims(bearer, &utils.GoTrueClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(api.config.CourierJWTSecret), nil
		})

		if err != nil {
			return utils.NewException(utils.ErrorStatusForbidden, fiber.StatusForbidden, err.Error())
		}

		claims, ok := token.Claims.(*utils.GoTrueClaims)

		if !ok {
			return utils.NewException(utils.ErrorInvalidToken, fiber.StatusNotAcceptable, "Your token is not valid")
		}

		ctx.Locals("claims", claims)
		ctx.Locals("token", token)
	} else {
		return utils.NewException(utils.ErrorStatusForbidden, fiber.StatusForbidden, "Only authorized requests are permitted")
	}

	return ctx.Next()
}
