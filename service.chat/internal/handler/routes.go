// Code generated by goctl. DO NOT EDIT.
package handler

import (
	"net/http"

	"service.chat/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/user/login",
				Handler: loginHandler(serverCtx),
			},
		},
	)

	server.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{serverCtx.Usercheck},
			[]rest.Route{
				{
					Method:  http.MethodPatch,
					Path:    "/pushtoken",
					Handler: updateTokenHandler(serverCtx),
				},
				{
					Method:  http.MethodGet,
					Path:    "/user/profile/:uid",
					Handler: getProfileHandler(serverCtx),
				},
				{
					Method:  http.MethodGet,
					Path:    "/messages/:chat_id",
					Handler: getMessageListHandler(serverCtx),
				},
				{
					Method:  http.MethodGet,
					Path:    "/conversations",
					Handler: getConversationListHandler(serverCtx),
				},
				{
					Method:  http.MethodPatch,
					Path:    "/conversation/:chat_id/last_read_time",
					Handler: updateLastReadTimeHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/conversation",
					Handler: newConversationHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/message",
					Handler: newChatMessageHandler(serverCtx),
				},
			}...,
		),
		rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
	)
}
