package tips

import (
	"context"

	"github.com/aymenworks/ProjectCookingTips-GoFromScratch/src/domain/entities"
)

type Service interface {
	GetAll(ctx context.Context) (entities.Tips, error)
	GetByID(ctx context.Context, id uint) (*entities.Tip, error)
}
