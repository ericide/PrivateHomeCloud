package logic

import (
	"context"

	"service.chat/internal/svc"
	"service.chat/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMessageListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetMessageListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMessageListLogic {
	return &GetMessageListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMessageListLogic) GetMessageList(req *types.MessagePageRequest) (resp *types.NormalListResponse, err error) {
	start := req.Size * req.Page
	list, err := l.svcCtx.ConversationMessageModel.Range(req.ChatId, start, req.Size)
	if err != nil {
		return nil, err
	}

	var list2 []types.RespConversationMessage
	for _, x := range list {
		list2 = append(list2, types.RespConversationMessage{
			Id:       x.Id,
			ChatId:   x.ChatId,
			Type:     x.Type,
			SenderId: x.SenderId,
			Content:  x.Content,
			SendTime: x.SendTime,
		})
	}

	return &types.NormalListResponse{
		Data: list2,
	}, nil
}
