package mysql

import (
	"errors"
	"fmt"

	"github.com/aymenworks/ProjectCookingTips-GoFromScratch/src/domain/entities"
	apperrors "github.com/aymenworks/ProjectCookingTips-GoFromScratch/src/errors"
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

func (r *MysqlTipsRepository) GetAll() (entities.Tips, error) {
	var tips entities.Tips
	result := r.db.Find(&tips)
	if result.Error != nil {
		return nil, apperrors.Stack(result.Error)
	}

	return tips, nil
}

func (r *MysqlTipsRepository) MustGetByID(id uint) (*entities.Tip, error) {
	var tip *entities.Tip
	if err := r.db.First(&tip, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.Wrap(apperrors.NotFound, fmt.Sprintf("id=%d", id))
		} else {
			return nil, apperrors.Wrap(err, fmt.Sprintf("id=%d", id))
		}
	}

	return tip, nil
}

func (r *MysqlTipsRepository) Create(tip *entities.Tip) error {
	if err := r.db.Create(&tip).Error; err != nil {
		return apperrors.Stack(err)
	}

	return nil
}

func (r *MysqlTipsRepository) DeleteByID(id uint) error {
	var tip entities.Tip
	nbRows := r.db.Delete(&tip, id).RowsAffected
	if nbRows == 0 {
		return apperrors.Wrap(apperrors.NotFound, fmt.Sprintf("id = %d", id))
	}

	return nil
}
