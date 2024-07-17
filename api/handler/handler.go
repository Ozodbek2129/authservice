package handler

import (
	"auth/storage/postgres"
	"database/sql"
	"github.com/go-redis/redis/v8"
)

type Handler struct {
	AuthUser *postgres.AuthUser
	Redis *redis.Client
}

func NewHanler(db *sql.DB,redisClient *redis.Client) *Handler{
	return &Handler{
		AuthUser: postgres.NewAuthUser(db),
		Redis: redisClient,
	}
}