package logic

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/edganiukov/apns"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"service.chat/internal/model"
	"time"

	"service.chat/internal/svc"
	"service.chat/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type NewChatMessageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewNewChatMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NewChatMessageLogic {
	return &NewChatMessageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *NewChatMessageLogic) NewChatMessage(req *types.SendMessageRequest) (resp *types.Response, err error) {
	senderUserId, ok := l.ctx.Value("user_id").(string)
	if !ok {
		return nil, errors.New("user id not exist")
	}

	plist, err := l.svcCtx.ConversationModel.QueryByChatId(l.ctx, req.ChatId)

	fmt.Println("list: ", plist)

	if err != nil {
		fmt.Println("list: zheli ")
		return nil, err
	}
	if len(*plist) == 0 {
		fmt.Println("errors: ", plist)
		return nil, errors.New("chat id is invalid")
	}

	cmItem := model.ConversationMessage{
		Id:         uuid.New().String(),
		ChatId:     req.ChatId,
		Type:       req.Type,
		SenderId:   senderUserId,
		Content:    req.Content,
		CreateTime: time.Now(),
	}

	l.svcCtx.ConversationMessageModel.Insert(l.ctx, &cmItem)

	pc := types.WSPushMessage{
		Type:              "MESSAGE",
		MessageChatId:     req.ChatId,
		MessageType:       cmItem.Type,
		MessageContent:    cmItem.Content,
		MessageSenderId:   cmItem.SenderId,
		MessageCreateTime: cmItem.CreateTime,
	}
	pushString, _ := json.Marshal(pc)
	for _, conItem := range *plist {
		l.SendMessage(conItem.OwnerId, "", pushString)
		if conItem.OwnerId != senderUserId {
			go l.SendPushToClient(pc, conItem.OwnerId)
		}
	}

	return &types.Response{}, nil
}

func (l *NewChatMessageLogic) SendMessage(userId string, jwtId string, content []byte) {
	l.svcCtx.ChannelChat <- &types.ChannelMessage{
		UserId:  userId,
		JwtId:   jwtId,
		Content: content,
	}
}

func (l *NewChatMessageLogic) SendPushToClient(message types.WSPushMessage, userId string) {
	user, err := l.svcCtx.UserModel.FindOne(l.ctx, userId)
	if err != nil {
		return
	}

	devices, err := l.svcCtx.UserLoginRecordModel.QueryByUser(userId)
	if err != nil {
		return
	}

	for _, item := range devices {
		if item.PushToken != "" {
			l.doPush(user.Name, message.MessageContent, item.PushToken)
		}
	}

}

func (l *NewChatMessageLogic) doPush(title string, content string, pushToken string) {
	fmt.Println("SendPushToClient")

	certificate, _ := tls.LoadX509KeyPair(l.svcCtx.Config.Push.CERT, l.svcCtx.Config.Push.KEY)
	//l.svcCtx.Config.Push.CERT
	c, err := apns.NewClient(
		apns.WithCertificate(certificate),
		apns.WithBundleID("com.cabital.cabital.debug.h5.container"),
		apns.WithMaxIdleConnections(10),
		apns.WithTimeout(5*time.Second),
		apns.WithEndpoint("https://api.sandbox.push.apple.com:443"),
	)
	if err != nil {
		fmt.Println(err)
	}
	badge := 100
	resp, err := c.Send(pushToken,
		apns.Payload{
			APS: apns.APS{
				Alert: apns.Alert{
					Title: title,
					Body:  content,
				},
				Badge: &badge,
			},
		},
		apns.WithExpiration(10),
		apns.WithCollapseID("test-collapse-id"),
		apns.WithPriority(5),
	)

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(resp)
}
