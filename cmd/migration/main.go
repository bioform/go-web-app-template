// This is custom goose binary with sqlite3 support only.

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/pressly/goose/v3"

	// Invoke init() functions within migrations pkg.
	_ "github.com/bioform/go-web-app-template/db/migrations"
	"github.com/bioform/go-web-app-template/pkg/database"
)

var (
	flags = flag.NewFlagSet("goose", flag.ExitOnError)
	dir   = flags.String("dir", "db/migrations", "directory with migration files")
)

func main() {
	err := flags.Parse(os.Args[1:])
	if err != nil {
		log.Fatalf("failed to parse flags: %v", err)
	}
	args := flags.Args()

	if len(args) < 1 {
		flags.Usage()
		return
	}

	ctx := context.Background()
	db, err := database.Get(ctx).DB()
	if err != nil {
		log.Fatalf("failed to open DB: %v", err)
	}

	command := args[0]

	arguments := []string{}
	if len(args) > 1 {
		arguments = append(arguments, args[1:]...)
	}

	fmt.Println("Command: ", command, "Arguments: ", arguments)

	if err := goose.RunContext(ctx, command, db, *dir, arguments...); err != nil {
		log.Fatalf("goose %v: %v", command, err)
	}
}
