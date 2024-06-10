package service

import (
	redis2 "car-parking-api/config/redis"
	"context"
	"encoding/json"
	"errors"
	"github.com/go-redis/redis/v8"
)

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UserService struct {
	redisClient redis2.Client
}

func NewUserService(redisClient redis2.Client) *UserService {
	return &UserService{redisClient: redisClient}
}

func (s *UserService) GetUserDetails(ctx context.Context, userID string) (*User, error) {
	val, err := s.redisClient.Get(ctx, userID).Result()
	if err == redis.Nil {
		return nil, errors.New("user not found")
	} else if err != nil {
		return nil, err
	}

	var user User
	if err := json.Unmarshal([]byte(val), &user); err != nil {
		return nil, err
	}

	return &user, nil
}
