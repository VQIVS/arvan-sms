package main

import (
	"flag"
	"log"
	"os"
	"sms-dispatcher/api/handler/http"
	"sms-dispatcher/app"
	"sms-dispatcher/config"
)

var configPath = flag.String("config", "config.json", "service configuration file")

func main() {
	flag.Parse()

	if v := os.Getenv("CONFIG_PATH"); len(v) > 0 {
		*configPath = v
	}

	c := config.MustReadConfig(*configPath)

	appContainer := app.NewMustApp(c)

	log.Fatal(http.Run(appContainer, c.Server))
}
