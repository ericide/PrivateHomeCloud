package logic

import (
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"path/filepath"
	"service.file/internal/svc"
	"service.file/internal/types"
	"time"
)

type GetFileListLogic struct {
	svcCtx *svc.ServiceContext
}

func NewGetFileListLogic( svcCtx *svc.ServiceContext) *GetFileListLogic {
	return &GetFileListLogic{
		svcCtx: svcCtx,
	}
}

func (l *GetFileListLogic) Do(context *gin.Context) (resp interface{}, err error) {

	log.Println(context.Params.ByName("path"))

	path := filepath.Join(l.svcCtx.Config.PhysicalPath, context.Params.ByName("path"))

	log.Println("path:", path)

	f, err := os.Open(path)
	if err != nil {
		log.Println(err)
		return nil, types.ErrNotFound
	}
	files, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		log.Println(err)
		return nil, types.ErrNotFound
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

	return &types.DataResponse{
		Data: rfiles,
	}, nil
}
