package user

import (
	"context"

	"github.com/aymenworks/ProjectCookingTips-GoFromScratch/src/domain/entities"
	"github.com/aymenworks/ProjectCookingTips-GoFromScratch/src/errors"
	"github.com/aymenworks/ProjectCookingTips-GoFromScratch/src/infra/caches"
	"github.com/aymenworks/ProjectCookingTips-GoFromScratch/src/repositories"
	"github.com/google/uuid"
)

type UserService struct {
	repository repositories.UserRepository
	cacheClt   caches.Cache
}

func NewUserService(repository repositories.UserRepository, cacheClt caches.Cache) Service {
	return &UserService{
		repository: repository,
		cacheClt:   cacheClt,
	}
}

func (s *UserService) Create(ctx context.Context, username, password string) (*entities.User, error) {
	u, err := s.repository.Create(ctx, uuid.NewString(), username, password)
	if err != nil {
		return nil, errors.Stack(err)
	}
	return u, nil
}

func (s *UserService) MustGetByUsername(ctx context.Context, username string) (*entities.User, error) {
	u, err := s.GetByUsername(ctx, username)
	if err != nil {
		return nil, errors.Stack(err)
	}

	if u == nil {
		return nil, errors.Stack(errors.NotFound)
	}

	return u, nil
}

func (s *UserService) GetByUsername(ctx context.Context, username string) (*entities.User, error) {
	u, err := s.repository.GetByUsername(ctx, username)
	if err != nil {
		return nil, errors.Stack(err)
	}
	return u, nil
}

func (s *UserService) VerifyAccessToken(ctx context.Context, uuid string) bool {
	key := "access_token:" + uuid
	t, err := s.cacheClt.Get(ctx, key)
	if err != nil {
		return false
	}
	if t == nil {
		return false
	}
	return true
}
