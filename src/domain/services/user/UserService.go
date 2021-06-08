package user

import (
	"context"

	"github.com/aymenworks/ProjectCookingTips-GoFromScratch/src/domain/entities"
	"github.com/aymenworks/ProjectCookingTips-GoFromScratch/src/errors"
	"github.com/aymenworks/ProjectCookingTips-GoFromScratch/src/repositories"
	"github.com/google/uuid"
)

type UserService struct {
	repository repositories.UserRepository
}

func NewUserService(repository repositories.UserRepository) Service {
	return &UserService{
		repository: repository,
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
