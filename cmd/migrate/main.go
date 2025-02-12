package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

const (
	dialect = "postgres"
)

var (
	flags = flag.NewFlagSet("migrate", flag.ExitOnError)
	dir   = flags.String("dir", "migrations", "directory to store migrations")
)

func main() {
	ctx := context.Background()

	flags.Usage = usage
	_ = flags.Parse(os.Args[1:])

	args := flags.Args()
	if len(args) == 0 || args[0] == "-h" || args[0] == "--help" {
		usage()
		return
	}

	command := args[0]

	dbConnect := os.Getenv("DB_CONNECT")
	if dbConnect == "" {
		dbConnect = "user=postgres password=postgres dbname=coffee_shop sslmode=disable host=localhost port=5433"
	}

	db, err := goose.OpenDBWithDriver(dialect, dbConnect)

	if err != nil {
		log.Fatal(err.Error())
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Fatal(err.Error())
		}
	}()

	if err := goose.RunContext(ctx, command, db, *dir, args[1:]...); err != nil {
		log.Fatalf("migrate %v: %v", command, err)
	}

}

func usage() {
	fmt.Println(usagePrefix)
	flags.PrintDefaults()
	fmt.Println(usageCommands)
}

var (
	usagePrefix = `Usage: migrate COMMAND
Examples:
    migrate status
`

	usageCommands = `
Commands:
    up                   Migrate the DB to the most recent version available
    up-by-one            Migrate the DB up by 1
    up-to VERSION        Migrate the DB to a specific VERSION
    down                 Roll back the version by 1
    down-to VERSION      Roll back to a specific VERSION
    redo                 Re-run the latest migration
    reset                Roll back all migrations
    status               Dump the migration status for the current DB
    version              Print the current version of the database
    create NAME [sql|go] Creates new migration file with the current timestamp
    fix                  Apply sequential ordering to migrations`
)
