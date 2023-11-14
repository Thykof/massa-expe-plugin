package main

import (
	"os"
	"os/signal"
	"syscall"

	"expe-plugin/api"
	"expe-plugin/api/server/restapi"
	"expe-plugin/api/server/restapi/operations"

	"github.com/bluele/gcache"
	"github.com/go-openapi/loads"
	"github.com/massalabs/station-massa-hello-world/pkg/plugin"
)

func initializeAPI() *restapi.Server {
	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "") // loads the swagger file
	if err != nil {
		panic(err)
	}

	gc := gcache.New(20).
		LRU().
		Build()

	pluginAPI := operations.NewExpePluginAPI(swaggerSpec) // initializes the API
	pluginAPI.ExpeHandler = api.NewExpe(gc)
	server := restapi.NewServer(pluginAPI) // creates the server

	server.ConfigureAPI() // configures the API

	return server
}

func main() {
	quit := make(chan bool)           // creates a channel to receive the interrupt signal
	intSig := make(chan os.Signal, 1) // notifies the channel when the interrupt signal is received
	signal.Notify(intSig, syscall.SIGINT, syscall.SIGTERM)

	server := initializeAPI() // initializes the API

	listener, err := server.HTTPListener()
	if err != nil {
		panic(err)
	}

	plugin.RegisterPlugin(listener) // registers the plugin

	if err := server.Serve(); err != nil {
		panic(err)
	}

	<-intSig
	quit <- true
}
