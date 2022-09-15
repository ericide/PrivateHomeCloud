package logic

import (
	"context"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"service.chat/internal/model"
	"strings"
	"time"

	"service.chat/internal/svc"
	"service.chat/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type NewConversationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewNewConversationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NewConversationLogic {
	return &NewConversationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *NewConversationLogic) NewConversation(req *types.CreateChannelRequest) (resp *types.Response, err error) {
	myUserId, ok := l.ctx.Value("user_id").(string)
	if !ok {
		return nil, errors.New("user id not exist")
	}

	chatId := l.get1To1ChatId(myUserId, req.UserId)

	list, err := l.svcCtx.ConversationModel.QueryByChatId(l.ctx, chatId)

	if err != nil {
		return nil, err
	}
	if len(*list) != 0 {
		return nil, nil
	}

	l.svcCtx.ConversationModel.Insert(l.ctx, &model.Conversation{
		Id:           uuid.New().String(),
		Type:         "1",
		ChatId:       chatId,
		OwnerId:      myUserId,
		OppoId:       req.UserId,
		Name:         "",
		LastReadTime: time.Now().UnixMicro(),
		CreateTime:   time.Now(),
	})

	l.svcCtx.ConversationModel.Insert(l.ctx, &model.Conversation{
		Id:           uuid.New().String(),
		Type:         "1",
		ChatId:       chatId,
		OwnerId:      req.UserId,
		OppoId:       myUserId,
		Name:         "",
		LastReadTime: time.Now().UnixMicro(),
		CreateTime:   time.Now(),
	})

	return &types.Response{}, nil
}

func (l *NewConversationLogic) get1To1ChatId(userId1 string, userId2 string) string {
	var str string = ""
	if strings.Compare(userId1, userId2) > 0 {
		str = userId2 + "_" + userId1
	} else {
		str = userId1 + "_" + userId2
	}

	return uuid.NewMD5(uuid.NameSpaceDNS, []byte(str)).String()
}
