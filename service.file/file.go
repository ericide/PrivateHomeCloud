package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"service.file/middleware"
	"service.file/middleware/types"
	"time"
)

func main() {
	router := gin.Default()

	router.Use(middleware.LoggerHandler())

	router.POST("file", func(context *gin.Context) {

	})

	router.GET("filelist/*path", func(context *gin.Context) {
		fmt.Println(context.Params.ByName("path"))

		path := filepath.Join(rootPath, context.Params.ByName("path"))

		var rfiles []types.RespFile

		f, err := os.Open(path)
		if err != nil {
			log.Fatal(err)
		}
		files, err := f.Readdir(-1)
		f.Close()
		if err != nil {
			log.Fatal(err)
		}

		for _, file := range files {
			//fmt.Println(file.Name())
			rfiles = append(rfiles, types.RespFile{
				Name:    file.Name(),
				Size:    file.Size(),
				ModTime: file.ModTime().Format(time.RFC3339),
				IsDir:   file.IsDir(),
			})
		}

		context.JSON(http.StatusOK, types.DataResponse{
			Data: rfiles,
		})
	})

	router.Static("file/", rootPath)

	router.Run(":8000")
}
