package middleware

import (
	"github.com/gofiber/fiber/v2"
)

func Middleware(ctx *fiber.Ctx) error {
	email := ctx.Get("Authorization")

	ctx.Locals("email", email)

	return ctx.Next()
}
