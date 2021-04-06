package mysql

import (
	"fmt"

	"github.com/aymenworks/ProjectCookingTips-GoFromScratch/src/domain/entities"
	"github.com/aymenworks/ProjectCookingTips-GoFromScratch/src/errors"
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
		return nil, errors.Wrap(result.Error, "error get all")
	}

	return tips, nil
}

func (r *MysqlTipsRepository) GetByID(id uint) (*entities.Tip, error) {
	var tip entities.Tip
	result := r.db.Find(&tip, id)
	if result.Error != nil {
		return nil, errors.Wrap(result.Error, "error get by id")
	}

	return &tip, nil
}

func (r *MysqlTipsRepository) Create(tip *entities.Tip) error {
	result := r.db.Create(&tip)
	if result.Error != nil {
		return errors.Wrap(result.Error, "error create tip")
	}

	return nil
}

func (r *MysqlTipsRepository) DeleteByID(id uint) error {
	var tip entities.Tip
	nbRows := r.db.Delete(&tip, id).RowsAffected
	if nbRows == 0 {
		return errors.Wrap(errors.NotFound, fmt.Sprintf("tip with id = %d not found", id))
	}

	return nil
}
