package cli

import (
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

	db := models.ConnectDatabase()
	db.DropTable(&models.User{})//TODO set dynamic inputs
	models.MigrateModels()

}
