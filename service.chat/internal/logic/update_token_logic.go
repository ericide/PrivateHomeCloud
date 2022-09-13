package logic

import (
	"context"
	"github.com/pkg/errors"
	"service.chat/internal/svc"
	"service.chat/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateTokenLogic {
	return &UpdateTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateTokenLogic) UpdateToken(req *types.UpdatePushTokenRequest) (resp *types.Response, err error) {
	jwtId, ok := l.ctx.Value("jwt_id").(string)
	if !ok {
		return nil, errors.New("jwt id not exist")
	}
	item, err := l.svcCtx.UserLoginRecordModel.FindOne(l.ctx, jwtId)
	if err != nil {
		return nil, err
	}

	item.PushToken = req.PushToken

	l.svcCtx.UserLoginRecordModel.Update(l.ctx, item)

	return
}
