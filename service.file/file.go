package main

import (
	"fmt"
	"github.com/hacdias/webdav/v4/lib"
	"go.uber.org/zap"
	"golang.org/x/net/webdav"
	"log"
	"net"
	"net/http"
	"os"
	"service.file/internal/config"
)

func main() {

	var cfg = &config.Config{
		Port:         8001,
		AccessToken:  os.Getenv("ACCESS_SECRET"),
		PhysicalPath: os.Getenv("FILE_ROOT_PATH"),
		Webdav: &webdav.Handler{
			Prefix: "/",
			FileSystem: lib.WebDavDir{
				Dir:     webdav.Dir("/Users/wu/Downloads/fusehfs-master"),
				NoSniff: false,
			},
			LockSystem: webdav.NewMemLS(),
		},
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		log.Fatal(err)
	}
	if err := http.Serve(listener, cfg); err != nil {
		log.Fatal("shutting server", zap.Error(err))
	}

}
