package main

import (
	"jwt-auth/cli"
	//"github.com/gin-gonic/gin"
)

func main() {

	cl := cli.CLI{}
	cl.Run()
}

