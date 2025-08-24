package auth

import (
	"fmt"

	//"strconv"
	//"strings"

	"github.com/gofiber/fiber/v2"
)

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	// Implement login logic
	fmt.Println("Login endpoint hit")
	return nil
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	// Implement registration logic
	fmt.Println("Register endpoint hit")
	return nil
}
