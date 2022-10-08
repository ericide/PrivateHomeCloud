package logic

import (
	"context"

	"service.chat/internal/svc"
	"service.chat/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetProfileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProfileLogic {
	return &GetProfileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetProfileLogic) GetProfile(req *types.GetProfileReq) (resp *types.RespUser, err error) {
	user, err := l.svcCtx.UserModel.FindOne(l.ctx, req.Uid)
	if err != nil {
		return nil, err
	}

	userModel := types.RespUser{
		Id:        user.Id,
		Name:      user.Name,
		AvatarUrl: user.AvatarUrl,
	}

	return &types.RespUser{
		Id:        userModel.Id,
		Name:      userModel.Name,
		AvatarUrl: userModel.AvatarUrl,
	}, nil
}
