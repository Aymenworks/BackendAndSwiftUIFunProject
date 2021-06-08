package mysql

import (
	"context"
	"errors"

	"github.com/aymenworks/ProjectCookingTips-GoFromScratch/src/domain/entities"
	apperrors "github.com/aymenworks/ProjectCookingTips-GoFromScratch/src/errors"
	"gorm.io/gorm"
)

type MysqlUserRepository struct {
	db *gorm.DB
}

func NewMysqlUserRepository(db *gorm.DB) *MysqlUserRepository {
	return &MysqlUserRepository{
		db: db,
	}
}

func (r *MysqlUserRepository) Create(ctx context.Context, uuid, username, password string) (*entities.User, error) {
	u := &entities.User{
		UUID:     uuid,
		Username: username,
		Password: password,
	}

	if err := r.db.Create(&u).Error; err != nil {
		return nil, apperrors.Stack(err)
	}

	return u, nil
}

func (r *MysqlUserRepository) GetByUsername(ctx context.Context, username string) (*entities.User, error) {
	var u *entities.User
	if err := r.db.
		Where("username = ?", username).
		First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, apperrors.Stack(err)
		}
	}

	return u, nil
}
