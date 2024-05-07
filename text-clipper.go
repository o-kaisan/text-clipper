package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/o-kaisan/text-clipper/common"
	"github.com/o-kaisan/text-clipper/text"
	"github.com/o-kaisan/text-clipper/tui"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	Version = "1.0.3"
)

func openSqlite() (*gorm.DB, error) {

	// デフォルトのデータベースパス
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("unable to find home directory: %w", err)
	}
	defaultDBPath := filepath.Join(homeDir, ".text-clipper", "text-clipper.db")

	// 環境変数で指定されたパスがあればそれを使用
	dbPath := os.Getenv("TEXT_CLIPPER_DB_PATH")
	if dbPath == "" {
		dbPath = defaultDBPath
	}

	// データベースディレクトリの作成
	dbDir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return nil, fmt.Errorf("unable to create database directory: %w", err)
	}

	// データベースに接続
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("unable to open database: %w", err)
	}

	err = db.AutoMigrate(&text.Text{})
	if err != nil {
		return db, fmt.Errorf("unable to migrate database: %w", err)
	}
	return db, nil
}

func main() {

	// Parsing options
	opts := common.ParseOptions(os.Args)

	// Setting log level (all logging must be after this line)
	common.SetupGlobalLogger(!opts.Debug)

	log.Printf("Original arguments: %#v", os.Args)
	log.Printf("Parsed options: %#v", opts) // must be after affecting -debug option

	switch {
	case opts.Help:
		fmt.Println(common.USAGE)
		os.Exit(0)
	case opts.Version:
		fmt.Println("text-clipper version " + Version)
		os.Exit(0)
	}

	db, err := openSqlite()
	if err != nil {
		log.Fatal(err)
	}
	tr := text.GormRepository{DB: db}
	tui.StartTea(tr)
}
