package cli

import (
	"fmt"
	"flag"
	"os"
	"log"
)

// CLI responsible for processing command line arguments
type CLI struct{}

func (cli *CLI) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  migrate -db FRESH - Fresh migrate")
	fmt.Println("  serve -port PORT - Port at input port")
}

func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}

func (cli *CLI) Run() {

	cli.validateArgs()

	getMigrateCmd := flag.NewFlagSet("migrate", flag.ExitOnError)
	getServeCmd := flag.NewFlagSet("serve", flag.ExitOnError)


	getMigrateCmdAction := getMigrateCmd.String("db", "", "the action for migrate operation")

	switch os.Args[1] {
	case "migrate":
		err := getMigrateCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "serve":
		err := getServeCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		cli.printUsage()
		os.Exit(1)
	}

	if getMigrateCmd.Parsed() {
		if *getMigrateCmdAction == "" {
			getMigrateCmd.Usage()
			os.Exit(1)
		}
		cli.Migrate(*getMigrateCmdAction)
	}

	if getServeCmd.Parsed() {
		cli.Serve()
	}




}

