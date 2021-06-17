package repositories

import (
	"context"

	"github.com/aymenworks/ProjectCookingTips-GoFromScratch/src/domain/entities"
)

type UserRepository interface {
	Create(ctx context.Context, uuid, username, password string) (*entities.User, error)
	GetByUsername(ctx context.Context, username string) (*entities.User, error)
	MustGetByUUID(ctx context.Context, uuid string) (*entities.User, error)
}
