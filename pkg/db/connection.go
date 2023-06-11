package db

import (
	"fmt"

	"github.com/nikhilnarayanan623/video-streaming-clean-arch/pkg/config"
	"github.com/nikhilnarayanan623/video-streaming-clean-arch/pkg/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDatabase(cfg *config.Config) (db *gorm.DB, err error) {

	dsn := fmt.Sprintf("host=%s user=%s dbname=%s port=%s password=%s", cfg.DBHost, cfg.DBUser, cfg.DBName, cfg.DBPort, cfg.DbPassword)

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database \nerror:%w", err)
	}

	err = db.AutoMigrate(
		&domain.Video{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to migrate tables \nerror:%w", err)
	}

	return
}
