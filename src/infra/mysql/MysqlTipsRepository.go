package mysql

import (
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

func (r *MysqlTipsRepository) GetAll() (entities.Tips, error) {
	var tips entities.Tips
	result := r.db.Find(&tips)
	if result.Error != nil {
		return nil, result.Error
	}

	return tips, nil
}

func (r *MysqlTipsRepository) GetByID(id uint) (*entities.Tip, error) {
	var tip entities.Tip
	result := r.db.Find(&tip, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return &tip, nil
}

func (r *MysqlTipsRepository) Create(tip *entities.Tip) error {
	result := r.db.Create(&tip)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *MysqlTipsRepository) DeleteByID(id uint) error {
	r.db.Delete(&entities.Tip{}, id)
	return nil // TODO: check if it's possible to know whether some deletion happened?
}
