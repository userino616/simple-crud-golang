package repository

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Config struct {
	Username string
	Password string
	DBName   string
}

func NewMysqlDB(cfg *Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@/%s", cfg.Username, cfg.Password, cfg.DBName)
	db, err := gorm.Open(mysql.Open(dsn))

	if err != nil {
		return nil, err
	}

	//if err := db.AutoMigrate(&models.User{}, &models.Post{}); err != nil {
	//	return nil, err
	//}

	return db, nil
}
