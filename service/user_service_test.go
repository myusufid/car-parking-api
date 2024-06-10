package service

import (
	"car-parking-api/config/redis"
	"context"
	"encoding/json"
	"github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetUserDetails(t *testing.T) {
	rdb, mock := redismock.NewClientMock()

	mockRedisClient := redis.Client{Client: rdb}

	userService := NewUserService(mockRedisClient)

	user := &User{
		ID:    "123",
		Name:  "John Doe",
		Email: "john.doe@example.com",
	}

	userJSON, err := json.Marshal(user)
	require.NoError(t, err)

	mock.ExpectGet("123").SetVal(string(userJSON))

	ctx := context.Background()
	result, err := userService.GetUserDetails(ctx, "123")
	require.NoError(t, err)
	assert.Equal(t, user, result)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestGetUserDetails_UserNotFound(t *testing.T) {
	rdb, mock := redismock.NewClientMock()

	mockRedisClient := redis.Client{Client: rdb}

	userService := NewUserService(mockRedisClient)

	mock.ExpectGet("123").RedisNil()

	ctx := context.Background()
	result, err := userService.GetUserDetails(ctx, "123")
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "user not found", err.Error())

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}
