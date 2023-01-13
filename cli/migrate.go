package cli

import (
	"jwt-auth/initializers"
	"os"
	"log"
	"jwt-auth/models"
)

func (cli *CLI) Migrate(action string){

	switch action {
	case "fresh":
		Fresh()
	default:
		log.Panic("invalid migrate cli argument")
		os.Exit(1)
	}

}


func Fresh(){

	config,err := initializers.LoadConfig(".")
	if err != nil{
		log.Panic(err)
	}

	initializers.ConnectDB(&config)
	initializers.DB.DropTable(&models.User{})//TODO set dynamic inputs
	initializers.MigrateModels()

}
