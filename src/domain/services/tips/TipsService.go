package tips

import (
	"context"

	"github.com/aymenworks/ProjectCookingTips-GoFromScratch/src/domain/entities"
	"github.com/aymenworks/ProjectCookingTips-GoFromScratch/src/errors"
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
	tips, err := s.repository.GetAll()
	if err != nil {
		return nil, errors.Wrap(err, "error get all")
	}
	return tips, nil
}

func (s *TipsService) MustGetByID(ctx context.Context, id uint) (*entities.Tip, error) {
	tip, err := s.repository.MustGetByID(id)
	if err != nil {
		return nil, errors.Wrap(err, "error get by id")
	}
	return tip, nil
}

func (s *TipsService) Create(ctx context.Context, name, path string) (*entities.Tip, error) {
	tip := &entities.Tip{
		Name:      name,
		ImagePath: path,
	}
	err := s.repository.Create(tip)
	if err != nil {
		return nil, errors.Wrap(err, "error create")
	}
	return tip, nil
}

func (s *TipsService) DeleteByID(ctx context.Context, id uint) error {
	err := s.repository.DeleteByID(id)
	if err != nil {
		return errors.Wrap(err, "error delete by id")
	}
	return nil
}
