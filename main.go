package main

import (
	"DoodleDropsBackend/initializers"
	"DoodleDropsBackend/routes"
	"fmt"
)

const (
	host = "0.0.0.0"
	port = 8000
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDb()
	initializers.SyncDatabase()
}

func main() {
	routes := routes.NewRoutes()

	routes.SetupRoutes()
	address := fmt.Sprintf("%s:%d", host, port)
	routes.Run(address)
}
