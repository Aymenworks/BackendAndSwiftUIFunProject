package tips

import (
	"context"

	"github.com/aymenworks/ProjectCookingTips-GoFromScratch/src/domain/entities"
)

type Service interface {
	GetAll(ctx context.Context) (entities.Tips, error)
	MustGetByID(ctx context.Context, id uint) (*entities.Tip, error)
	Create(ctx context.Context, name string) (*entities.Tip, error)
	DeleteByID(ctx context.Context, id uint) error
}
