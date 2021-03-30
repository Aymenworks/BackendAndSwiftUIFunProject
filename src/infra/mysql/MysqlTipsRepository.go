package mysql

import (
	"context"

	"github.com/aymenworks/ProjectCookingTips-GoFromScratch/src/domain/entities"
	"gorm.io/gorm"
)

type MysqlTipsRepository struct {
	db *gorm.DB
}

func NewMysqlTipsRepository(db *gorm.DB) *MysqlTipsRepository {
	return &MysqlTipsRepository{
		db: db,
	}
}

func (r *MysqlTipsRepository) GetAll(ctx context.Context) (entities.Tips, error) {
	var tips entities.Tips
	result := r.db.Find(&tips)
	if result.Error != nil {
		return nil, result.Error
	}

	return tips, nil
}

func (r *MysqlTipsRepository) GetByID(ctx context.Context, id uint) (*entities.Tip, error) {
	var tip entities.Tip
	result := r.db.Find(&tip, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}

	return &tip, nil
}
