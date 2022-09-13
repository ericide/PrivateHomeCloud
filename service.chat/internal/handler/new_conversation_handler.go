package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"service.chat/internal/logic"
	"service.chat/internal/svc"
	"service.chat/internal/types"
)

func newConversationHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateChannelRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewNewConversationLogic(r.Context(), svcCtx)
		resp, err := l.NewConversation(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
