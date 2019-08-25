package main

import (
	"flag"
	"fmt"
	"inmemoryStorage/cmd/storage/app"
	"inmemoryStorage/config"
)

var (
	port     = flag.String("port", ":2094", "App port")
	ttlCheck = flag.Int("ttl_check", 5, "TTL check interval")
)

func main() {
	flag.Parse()

	cfg := config.New(*ttlCheck, *port)

	application, err := app.New(cfg)
	if err != nil {
		return
	}

	if err := application.Run(); err != nil {
		fmt.Println("close application")
		return
	}

}
