package AdminMiddlewareAccess

import (
	"github.com/gofiber/fiber/v2"

	"github.com/dgrijalva/jwt-go"
	/* "net/http" */

)

const SecretKey = "secret"
func AdminMiddleware(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")

	token, err := jwt.ParseWithClaims(cookie, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "unauthenticated"})
	}

	claims := token.Claims.(*jwt.MapClaims)

	if (*claims)["IsAdmin"] != true {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "access denied, admin only"})
	}

	return c.Next()
}
