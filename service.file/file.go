package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"service.file/internal/config"
	"service.file/internal/handler"
	"service.file/internal/svc"
)

func main() {
	engine := gin.Default()

	var c = config.Config{
		Port: 8001,
		AccessToken:  os.Getenv("ACCESS_SECRET"),
		PhysicalPath: os.Getenv("FILE_ROOT_PATH"),
	}
	ctx := svc.NewServiceContext(c)

	handler.RegisterHandlers(engine, ctx)

	engine.Run( fmt.Sprintf(":%d", c.Port))
	log.Printf("server run : %d", c.Port)
}
