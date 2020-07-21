package main

import (
	"be_auth_go/app/controllers"
)

func main() {
	app := controllers.App{}
	app.RunServer()
}