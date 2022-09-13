package svc

import (
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/rest"
	"service.chat/internal/config"
	"service.chat/internal/middleware"
	"service.chat/internal/model"
	"service.chat/internal/types"
)

type ServiceContext struct {
	Config                   config.Config
	Usercheck                rest.Middleware
	UserModel                model.UserModel
	UserLoginRecordModel     model.UserLoginRecordModel
	ConversationModel        model.ConversationModel
	ConversationMessageModel model.ConversationMessageModel
	ChannelChat              chan *types.ChannelMessage
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.DataSource)

	return &ServiceContext{
		Config:                   c,
		Usercheck:                middleware.NewUsercheckMiddleware(model.NewUserModel(conn), model.NewUserLoginRecordModel(conn)).Handle,
		UserModel:                model.NewUserModel(conn),
		UserLoginRecordModel:     model.NewUserLoginRecordModel(conn),
		ConversationModel:        model.NewConversationModel(conn),
		ConversationMessageModel: model.NewConversationMessageModel(conn),
		ChannelChat:              make(chan *types.ChannelMessage),
	}
}
