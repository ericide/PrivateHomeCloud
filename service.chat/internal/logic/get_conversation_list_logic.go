package logic

import (
	"context"

	"service.chat/internal/svc"
	"service.chat/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetConversationListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetConversationListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetConversationListLogic {
	return &GetConversationListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetConversationListLogic) GetConversationList(req *types.NullRequest) (resp *types.NormalListResponse, err error) {
	myUserId, ok := l.ctx.Value("user_id").(string)
	if !ok {
		return nil, err
	}

	list, err := l.svcCtx.ConversationModel.QueryByUserId(l.ctx, myUserId)
	if err != nil {
		return nil, err
	}

	var list2 []types.RespConversation
	for _, x := range *list {

		oppoUser, err := l.svcCtx.UserModel.FindOne(l.ctx, x.OppoId)
		if err != nil {
			return nil, err
		}

		list2 = append(list2, types.RespConversation{
			Id:      x.Id,
			Type:    x.Type,
			ChatId:  x.ChatId,
			OwnerId: x.OwnerId,
			OppoUser: types.RespUser{
				Id:        oppoUser.Id,
				Name:      oppoUser.Name,
				AvatarUrl: oppoUser.AvatarUrl,
			},
			Name:         x.Name,
			LastReadTime: x.LastReadTime.String(),
			CreateTime:   x.CreateTime.String(),
		})
	}

	return &types.NormalListResponse{
		Data: list2,
	}, nil
}
