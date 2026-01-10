package api

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	"nudgepay/internal/auth"
)

func authRequired(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return fiber.NewError(fiber.StatusUnauthorized, "missing authorization header")
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid authorization header")
		}
		claims, err := auth.ParseToken(secret, parts[1])
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid token")
		}
		c.Locals("user_id", claims.UserID)
		c.Locals("org_id", claims.OrgID)
		return c.Next()
	}
}

func jsonErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	msg := "internal error"
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		msg = e.Message
	}
	return c.Status(code).JSON(fiber.Map{"error": msg})
}

func orgIDFrom(c *fiber.Ctx) string {
	v := c.Locals("org_id")
	if v == nil {
		return ""
	}
	return v.(string)
}

func userIDFrom(c *fiber.Ctx) string {
	v := c.Locals("user_id")
	if v == nil {
		return ""
	}
	return v.(string)
}
