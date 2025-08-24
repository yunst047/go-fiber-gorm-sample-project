package jwt

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"go-fiber-gorm-sample/config"
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(config.AccessTokenConfig.Secret)

type Claims struct {
	UserID    uint   `json:"user_id"`
	Role      string `json:"role"`
	TokenType string `json:"token_type"`
	jwt.RegisteredClaims
}

func GenerateToken(userID uint, role string) (string, string, error) {
	accessTokenExpiration := time.Now().UTC().Add(time.Duration(config.AccessTokenConfig.AccessTokenExpiration) * time.Second)
	refreshTokenExpiration := time.Now().UTC().Add(time.Duration(config.AccessTokenConfig.RefreshTokenExpiration) * time.Second)

	accessTokenClaims := &Claims{
		UserID:    userID,
		Role:      role,
		TokenType: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessTokenExpiration),
		},
	}

	refreshTokenClaims := &Claims{
		UserID:    userID,
		Role:      role,
		TokenType: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshTokenExpiration),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)

	accessTokenString, err := accessToken.SignedString(jwtSecret)
	if err != nil {
		log.Printf("Error signing access token: %v", err)
		return "", "", err
	}

	refreshTokenString, err := refreshToken.SignedString(jwtSecret)
	if err != nil {
		log.Printf("Error signing refresh token: %v", err)
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}

func ParseToken(tokenString string) (*Claims, error) {

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecret, nil
	})

	if err != nil {
		log.Printf("Error parsing token: %v", err)
		return nil, err
	}
	if !token.Valid {
		log.Printf("Invalid token: token is not valid")
		return nil, errors.New("invalid token")
	}

	if claims.UserID == 0 || claims.Role == "" || claims.ExpiresAt == nil {
		log.Printf("Error: claims are not properly set or missing required fields")
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}

func ValidateToken(tokenStr string) (*Claims, error) {
	return ParseToken(tokenStr)
}

func IsTokenExpired(claims *Claims) bool {
	if claims.ExpiresAt == nil {

		return true
	}
	log.Printf("Current time: %s", time.Now().UTC().String())
	log.Printf("Token expiry time: %s", claims.ExpiresAt.Time.String())
	return time.Until(claims.ExpiresAt.Time) <= 0
}

func RefreshAccessToken(refreshTokenString string) (string, error) {

	log.Printf("Received refresh token: %s", refreshTokenString)
	claims, err := ValidateToken(refreshTokenString)
	if err != nil {
		log.Printf("Error validating refresh token: %v", err)
		return "", errors.New("authentication failed")
	}

	if IsTokenExpired(claims) {
		log.Printf("Refresh token is expired")
		return "", errors.New("authentication failed")
	}

	newAccessToken, _, err := GenerateToken(claims.UserID, claims.Role)
	if err != nil {
		log.Printf("Error generating new access token: %v", err)
		return "", err
	}

	return newAccessToken, nil
}

func GetJWTClaims(c *fiber.Ctx) (*Claims, error) {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		log.Println("Authorization header is missing")
		return nil, errors.New("authentication failed")
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		log.Println("Bearer token is missing")
		return nil, errors.New("authentication failed")
	}

	log.Printf("Received token: %s", tokenString)

	claims, err := ParseToken(tokenString)
	if err != nil {
		log.Printf("Error parsing token: %v", err)

		log.Printf("Token string: %s", tokenString)
		return nil, errors.New("authentication failed")
	}

	log.Printf("Parsed claims: %+v", claims)
	return claims, nil
}

func JWTProtected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Bypass JWT authentication for OPTIONS (preflight) requests
		if c.Method() == fiber.MethodOptions {
			return c.Next() // Skip JWT validation for preflight requests
		}

		authHeader := c.Get("Authorization")
		token := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := ParseToken(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
		}

		if IsTokenExpired(claims) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token expired"})
		}

		c.Locals("userID", claims.UserID)
		return c.Next()
	}
}

func GenerateSecurePassword(length int) (string, error) {

	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(bytes)[:length], nil
}
