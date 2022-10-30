package main

import (
	"file/config"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
)

func main() {

	var cfg = &config.Config{
		Port:        50960,
		AccessToken: os.Getenv("ACCESS_SECRET"),
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		log.Fatal(err)
	}
	if err := http.Serve(listener, cfg); err != nil {
		log.Fatal("shutting server")
	}

}
