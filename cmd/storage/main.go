package main

import (
	"flag"
	"fmt"
	"inmemoryStorage/cmd/storage/app"
	"inmemoryStorage/config"
)

var (
	port         = flag.String("port", ":2094", "App port")
	ttlCheck     = flag.Int("ttl_check", 5, "TTL check interval")
	dumpFile     = flag.String("dump", "dump.db", "Dump file")
	dumpInterval = flag.Int("dump_interval", 5, "Dump interval")
)

func main() {
	flag.Parse()

	cfg := config.New(*ttlCheck, *port, *dumpFile, *dumpInterval)

	application, err := app.New(cfg)
	if err != nil {
		return
	}

	if err := application.Run(); err != nil {
		fmt.Println("close application")
		return
	}

}
