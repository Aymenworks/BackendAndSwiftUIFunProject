package repositories

import (
	"github.com/aymenworks/ProjectCookingTips-GoFromScratch/src/domain/entities"
)

type TipsRepository interface {
	Create(tip *entities.Tip) error
	DeleteByID(id uint) error
	GetAll() (entities.Tips, error)
	MustGetByID(id uint) (*entities.Tip, error)
}
