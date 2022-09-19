package logic

import (
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"path/filepath"
	"service.file/internal/svc"
	"service.file/internal/types"
)

type NewDirectoryLogic struct {
	svcCtx *svc.ServiceContext
}

func NewNewDirectoryLogic( svcCtx *svc.ServiceContext) *GetFileListLogic {
	return &GetFileListLogic{
		svcCtx: svcCtx,
	}
}

func (l *NewDirectoryLogic) Do(context *gin.Context) (resp interface{}, err error) {

	basePath := context.Request.FormValue("base_path")
	fileName := context.Request.FormValue("file_name")

	finalPath := filepath.Join(l.svcCtx.Config.PhysicalPath, basePath, fileName)

	log.Println(l.svcCtx.Config.PhysicalPath, basePath, fileName, finalPath)

	err = l.createDir(finalPath)

	if err != nil {
		return nil, err
	}

	return &types.DataResponse{
		Data: fileName,
	}, nil
}

func (l *NewDirectoryLogic) createDir(path string)  error {
	_exist, _err := l.hasDir(path)
	if _err != nil {
		log.Printf("获取文件夹异常 -> %v\n", _err)
		return _err
	}
	if _exist {
		log.Println("文件夹已存在！")
		return errors.New("文件夹已存在")
	} else {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			log.Printf("创建目录异常 -> %v\n", err)
			return err
		} else {
			return nil
		}
	}
}

func (l *NewDirectoryLogic) hasDir(path string) (bool, error) {
	_, _err := os.Stat(path)
	if _err == nil {
		return true, nil
	}
	if os.IsNotExist(_err) {
		return false, nil
	}
	return false, _err
}
