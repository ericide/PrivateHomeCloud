package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"service.chat/internal/logic"
	"service.chat/internal/svc"
	"service.chat/internal/types"
)

func updateLastReadTimeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdateConversationLastReadTimeRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewUpdateLastReadTimeLogic(r.Context(), svcCtx)
		resp, err := l.UpdateLastReadTime(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
