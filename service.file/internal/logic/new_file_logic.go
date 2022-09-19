package logic

import (
	"github.com/gin-gonic/gin"
	"log"
	"path/filepath"
	"service.file/internal/svc"
	"service.file/internal/types"
)

type NewFileLogic struct {
	svcCtx *svc.ServiceContext
}

func NewNewFileLogic(svcCtx *svc.ServiceContext) *NewFileLogic {
	return &NewFileLogic{
		svcCtx: svcCtx,
	}
}

func (l *NewFileLogic) Do(context *gin.Context) (resp interface{}, err error) {
	log.Printf("NewFileLogic\n")
	file, err := context.FormFile("file")
	if err != nil {
		return nil, err
	}
	// c.JSON(200, gin.H{"message": file.Header.Context})

	basePath := context.Request.FormValue("base_path")

	finalPath := filepath.Join(l.svcCtx.Config.PhysicalPath, basePath, file.Filename)

	log.Println(l.svcCtx.Config.PhysicalPath, basePath, file.Filename, finalPath)

	err = context.SaveUploadedFile(file, finalPath)

	if err != nil {
		log.Printf("err:", err)
		return nil, err
	}

	return &types.DataResponse{
		Data: file.Filename,
	}, nil
}
