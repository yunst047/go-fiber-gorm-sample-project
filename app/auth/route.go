package auth

import (
	"go-fiber-gorm-sample/jwt"

	"github.com/gofiber/fiber/v2"
)

func RegisterAuthRoutes(router fiber.Router, storage *AuthStorage) {
	protected := jwt.JWTProtected()
	auth := router.Group("/auth")
	auth.Post("/login", protected, LoginFunc(storage))
	auth.Post("/register", RegisterFunc(storage))

	// Add more routes as needed
}

func LoginFunc(storage *AuthStorage) fiber.Handler {
	handler := NewAuthHandler(storage)
	return handler.Login
}

func RegisterFunc(storage *AuthStorage) fiber.Handler {
	handler := NewAuthHandler(storage)
	return handler.Register
}
