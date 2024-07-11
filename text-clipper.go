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
	Version = "1.0.7"
)

func openSqlite() (*gorm.DB, error) {

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
		os.Exit(1)
	}
	tr := text.GormRepository{DB: db}
	tui.StartTea(tr)
}
