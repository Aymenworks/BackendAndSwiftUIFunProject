package tips

import (
	"context"

	"github.com/aymenworks/ProjectCookingTips-GoFromScratch/src/domain/entities"
	"github.com/aymenworks/ProjectCookingTips-GoFromScratch/src/repositories"
)

type TipsService struct {
	repository repositories.TipsRepository
}

func NewTipsService(repository repositories.TipsRepository) Service {
	return &TipsService{
		repository: repository,
	}
}

func (s *TipsService) GetAll(ctx context.Context) (entities.Tips, error) {
	return s.repository.GetAll(ctx)
	// TODO: If errors, look if we get the full stack
	// TODO: Implement stack trace error
}

func (s *TipsService) GetByID(ctx context.Context, id uint) (*entities.Tip, error) {
	return s.repository.GetByID(ctx, id)
}
