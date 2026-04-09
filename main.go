package main

import _ "github.com/lib/pq"

import (
	// "fmt"
	"log"
	"os"
	"database/sql"

	"github.com/MarunDArbaumont/blog-aggregator/internal/config"
	"github.com/MarunDArbaumont/blog-aggregator/internal/database"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	// fmt.Printf("Read config: %+v\n", cfg)

	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		log.Fatalf("error connecting to database: %v", err)
	}
	// fmt.Printf("Connected to: %v\n", db)
	dbQueries := database.New(db)

	mainState := &state{
		db: dbQueries,
		cfg: cfg,
	}	

	mainCommands := commands{
		command: make(map[string]func(*state, command) error),
	}

	mainCommands.register("login", handlerLogin)
	mainCommands.register("register", handlerRegister)
	mainCommands.register("reset", handlerReset)
	mainCommands.register("users", handlerUsers)
	mainCommands.register("agg", handlerAgg)
	mainCommands.register("addfeed", handlerAddFeed)
	mainCommands.register("feeds", handlerFeeds)
	mainCommands.register("follow", handlerFollow)
	mainCommands.register("following", handlerFollowing)

	if len(os.Args) < 2 {
		log.Fatalf("not enough arguments")
	}

	commandName := os.Args[1]
	commandArgs := os.Args[2:]

	mainCommand := command{
		name: commandName,
		args: commandArgs,
	}

	err = mainCommands.run(mainState, mainCommand)
	if err != nil {
		log.Fatal(err)
	}
}