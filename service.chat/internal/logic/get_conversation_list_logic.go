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

		oppoUserWrap := types.RespUser{
			Id:        oppoUser.Id,
			Name:      oppoUser.Name,
			AvatarUrl: oppoUser.AvatarUrl,
		}

		unreadCount, err := l.svcCtx.ConversationMessageModel.CountAfterTime(l.ctx, x.ChatId, x.LastReadTime)
		if err != nil {
			return nil, err
		}

		lastMessage, _ := l.svcCtx.ConversationMessageModel.LastMessage(l.ctx, x.ChatId)
		var lastMessageWrap *types.RespConversationMessage = nil
		if lastMessage != nil {
			lastMessageWrap = &types.RespConversationMessage{
				Id:         lastMessage.Id,
				ChatId:     lastMessage.ChatId,
				Type:       lastMessage.Type,
				SenderId:   lastMessage.SenderId,
				Content:    lastMessage.Content,
				CreateTime: lastMessage.CreateTime.String(),
			}
		}

		list2 = append(list2, types.RespConversation{
			Id:           x.Id,
			Type:         x.Type,
			ChatId:       x.ChatId,
			OwnerId:      x.OwnerId,
			OppoUser:     oppoUserWrap,
			LastMessage:  lastMessageWrap,
			UnreadCount:  *unreadCount,
			Name:         x.Name,
			LastReadTime: x.LastReadTime.String(),
			CreateTime:   x.CreateTime.String(),
		})
	}

	return &types.NormalListResponse{
		Data: list2,
	}, nil
}
