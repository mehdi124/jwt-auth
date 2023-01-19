package cli

import (
	"jwt-auth/utils/email"
	Redis "jwt-auth/utils/redis"
)

func (cli *CLI) Work(){

	work()

}

func work(){


	rdq := Redis.NewClient()
	//go email.StartProcessingEmails(rdq)
	email.StartProcessingEmails(rdq)
}