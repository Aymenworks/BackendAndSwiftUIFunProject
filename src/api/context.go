package api

import (
	"context"

	"go.uber.org/zap"
)

type ContextKey string

const (
	userUUIDKey        ContextKey = "tips_user_uuid_key"
	accessTokenUUIDKey ContextKey = "tips_access_token_uuid_key"
)

// TODO: check whether the address of the key is really unique with a custom type

func WithUserUUID(ctx context.Context, uuid string) context.Context {
	return context.WithValue(ctx, userUUIDKey, uuid)
}

func UserUUID(ctx context.Context) string {
	uuid, ok := ctx.Value(userUUIDKey).(string)
	handleUnsetValue(ok, "UserUUID")
	return uuid
}

func WithAccessTokenUUID(ctx context.Context, uuid string) context.Context {
	return context.WithValue(ctx, accessTokenUUIDKey, uuid)
}

func AccessTokenUUID(ctx context.Context) string {
	uuid, ok := ctx.Value(accessTokenUUIDKey).(string)
	handleUnsetValue(ok, "AccessTokenUUID")
	return uuid
}

func handleUnsetValue(ok bool, name string) {
	if !ok {
		zap.S().Panicf("fatal context error: %s not set", name)
	}
}
