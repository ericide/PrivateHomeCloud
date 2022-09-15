package logic

import (
	"context"
	"encoding/json"
	"errors"
	"service.chat/internal/defines"
	"service.chat/internal/model"
	"time"

	"service.chat/internal/svc"
	"service.chat/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateLastReadTimeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateLastReadTimeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLastReadTimeLogic {
	return &UpdateLastReadTimeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateLastReadTimeLogic) UpdateLastReadTime(req *types.UpdateConversationLastReadTimeRequest) (resp *types.Response, err error) {
	myUserId, ok := l.ctx.Value("user_id").(string)
	if !ok {
		return nil, errors.New("user id not exist")
	}

	cst, err := l.svcCtx.ConversationModel.QueryByChatIdUserId(l.ctx, req.ChatId, myUserId)

	if err != nil {
		return nil, err
	}

	cst.LastReadTime = time.Now().UnixMicro()

	l.svcCtx.ConversationModel.Update(l.ctx, cst)

	l.SendWSMessage(cst)

	return &types.Response{
		Message: "",
	}, nil
}

func (l *UpdateLastReadTimeLogic) SendWSMessage(cstItem *model.Conversation) {

	plist, err := l.svcCtx.ConversationModel.QueryByChatId(l.ctx, cstItem.ChatId)
	if err != nil {
		return
	}

	pc := types.WSUpdateReadTime{
		WSPushBase: types.WSPushBase{
			Type:    defines.WSType_Message,
			SubType: defines.WSSubType_UpdateTime,
		},
		ChatId:   cstItem.ChatId,
		UserId:   cstItem.OwnerId,
		ReadTime: cstItem.LastReadTime,
	}

	pushString, _ := json.Marshal(pc)
	for _, person := range *plist {
		l.svcCtx.ChannelChat <- &types.ChannelMessage{
			UserId:  person.OwnerId,
			JwtId:   "",
			Content: pushString,
		}
	}
}
