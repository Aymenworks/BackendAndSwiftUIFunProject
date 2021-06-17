package user

import (
	"context"

	"github.com/aymenworks/ProjectCookingTips-GoFromScratch/src/domain/entities"
)

type Service interface {
	Create(ctx context.Context, username, password string) (*entities.User, error)
	VerifyAccessToken(ctx context.Context, uuid string) bool
	GetByUsername(ctx context.Context, username string) (*entities.User, error)
	MustGetByUsername(ctx context.Context, username string) (*entities.User, error)
	MustGetByUUID(ctx context.Context, uuid string) (*entities.User, error)
}
