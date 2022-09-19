package handler

import (
	"github.com/gin-gonic/gin"
	"service.file/internal/logic"
	"service.file/internal/middleware"
	"service.file/internal/svc"
)

func RegisterHandlers(server *gin.Engine , serverCtx *svc.ServiceContext) {
	server.Use(middleware.LoggerHandler())

	server.GET("filelist/*path", createNormalHandler(logic.NewGetFileListLogic(serverCtx)) )

	server.POST("directory", createNormalHandler(logic.NewNewDirectoryLogic(serverCtx)))

	server.POST("file", createNormalHandler(logic.NewNewFileLogic(serverCtx)))

	server.Static("file/", serverCtx.Config.PhysicalPath)
}

type StandardHandler interface {
	Do(context *gin.Context) (interface{}, error)
}

