package user

import (
	"context"

	"github.com/aymenworks/ProjectCookingTips-GoFromScratch/src/domain/entities"
)

type Service interface {
	Create(ctx context.Context, username, password string) (*entities.User, error)
	GetByUsername(ctx context.Context, username string) (*entities.User, error)
	MustGetByUsername(ctx context.Context, username string) (*entities.User, error)
}
