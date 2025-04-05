package main

import (
	"fmt"
	"log"
	"os"

	"github.com/o-kaisan/text-clipper/common"
	"github.com/o-kaisan/text-clipper/di"
	"github.com/o-kaisan/text-clipper/infrastructure/sqlite"
	app "github.com/o-kaisan/text-clipper/interface/bubbletea"
)

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
		fmt.Println("text-clipper version " + common.GetVersionFromGit())
		os.Exit(0)
	}

	db, err := sqlite.NewSqlite()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	container := di.NewContainer(db)
	app.StartTea(container)
}
