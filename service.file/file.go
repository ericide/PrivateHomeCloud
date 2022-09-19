package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"service.file/internal/middleware"
	"service.file/internal/types"
	"time"
)

func main() {
	router := gin.Default()

	rootPath := os.Getenv("FILE_ROOT_PATH")

	router.Use(middleware.LoggerHandler())

	router.POST("file", func(context *gin.Context) {
		file, err := context.FormFile("file")
		if err != nil {
			context.JSON(500, gin.H{
				"error": "失败",
			})
			return
		}
		// c.JSON(200, gin.H{"message": file.Header.Context})

		basePath := context.Request.FormValue("base_path")

		finalPath := filepath.Join(rootPath, basePath, file.Filename)

		fmt.Println(rootPath, basePath, file.Filename, finalPath)

		err = context.SaveUploadedFile(file, finalPath)

		if err != nil {
			fmt.Println(err)
		}

		context.String(http.StatusOK, file.Filename)
	})

	router.GET("filelist/*path", func(context *gin.Context) {
		fmt.Println(context.Params.ByName("path"))

		path := filepath.Join(rootPath, context.Params.ByName("path"))

		fmt.Println("path:", path)

		f, err := os.Open(path)
		if err != nil {
			log.Println(err)
			context.JSON(http.StatusNotFound, gin.H{
				"error": "not found",
			})
			return
		}
		files, err := f.Readdir(-1)
		f.Close()
		if err != nil {
			log.Println(err)
			context.JSON(http.StatusNotFound, gin.H{
				"error": "not found",
			})
			return
		}

		var rfiles = []types.RespFile{}
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

	router.Run(":8001")
}
