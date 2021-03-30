package databases

import (
	"fmt"

	"github.com/aymenworks/ProjectCookingTips-GoFromScratch/src/domain/entities"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MySQLDatabase struct {
}

func NewMySQLDatabase(port, database string) *gorm.DB {
	dsn := fmt.Sprintf("root@tcp(127.0.0.1:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local", port, database)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		zap.S().Panicf("Error opening db %v\n", err)
	}

	zap.S().Debug("Db is ok")

	db.AutoMigrate(&entities.Tip{})

	return db
}
