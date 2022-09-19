package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"service.file/internal/types"
)

func createNormalHandler(handler StandardHandler ) gin.HandlerFunc {
	return func(context *gin.Context) {
		info, err := handler.Do(context)

		if err == types.ErrNotFound {
			context.JSON(http.StatusNotFound, gin.H{
				"message": err,
			})
			return
		}

		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{
				"message": err,
			})
			return
		}

		context.JSON(http.StatusOK, info)
	}
}
