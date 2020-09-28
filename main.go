package main

import (
	"flag"
	"job/presentation/controller"
	"job/presentation/core/config"
	"job/presentation/core/routes"
	"log"
	"net/http"
)

var (
	configPath = flag.String("config", "config.toml", "config file path")
)

func main() {
	conf, err := config.Read(*configPath)
	if err != nil {
		log.Println(err)
		return
	}

	env, err := controller.NewEnvironment()
	if err != nil {
		log.Println(err)
		return
	}

	router, err := routes.NewRouter(env, conf)
	if err != nil {
		log.Println(err)
		return
	}

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	log.Println("Listen and serve")
	server.ListenAndServe()
}
