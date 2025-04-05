package sqlite

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/o-kaisan/text-clipper/common"
	"github.com/o-kaisan/text-clipper/domain/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewSqlite() (*gorm.DB, error) {

	dbPath, err := common.GetPathFromPath("text-clipper.db")
	if err != nil {
		return nil, fmt.Errorf("failed in getting db path: %w", err)
	}
	dbDir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return nil, fmt.Errorf("unable to create database directory: %w", err)
	}

	// データベースに接続
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("unable to open database: %w", err)
	}

	err = db.AutoMigrate(&model.Clip{})
	if err != nil {
		return db, fmt.Errorf("unable to migrate database: %w", err)
	}
	return db, nil
}
