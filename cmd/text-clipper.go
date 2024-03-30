package main

import (
	"fmt"
	"log"
	"os"

	"github.com/o-kaisan/text-clipper/common"
	"github.com/o-kaisan/text-clipper/text"
	"github.com/o-kaisan/text-clipper/tui"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	Version = "1.0.0"
)

func openSqlite() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("text-clipper.db"), &gorm.Config{})
	if err != nil {
		return db, fmt.Errorf("unable to open database: %w", err)
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
