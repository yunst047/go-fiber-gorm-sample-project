package auth

import (
	"gorm.io/gorm"
)

type AuthStorage struct {
	db *gorm.DB
}

func NewAuthStorage(db *gorm.DB) *AuthStorage {
	return &AuthStorage{db: db}
}

type AuthHandler struct {
	storage *AuthStorage
}

func NewAuthHandler(storage *AuthStorage) *AuthHandler {
	return &AuthHandler{storage: storage}
}
