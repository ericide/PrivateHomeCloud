package logic

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/edganiukov/apns"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"service.chat/internal/defines"
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
	//fmt.Println("list: ", plist)

	if err != nil {
		return nil, err
	}
	if len(*plist) == 0 {
		fmt.Println("errors: ", plist)
		return nil, errors.New("chat id is invalid")
	}

	cmItem := model.ConversationMessage{
		Id:       uuid.New().String(),
		ChatId:   req.ChatId,
		Type:     req.Type,
		SenderId: senderUserId,
		Content:  req.Content,
		SendTime: time.Now().UnixMicro(),
	}

	l.svcCtx.ConversationMessageModel.Insert(l.ctx, &cmItem)

	l.SendWSMessage(plist, req.MessageClientId, cmItem)
	go l.SendPushNotificationToClient(plist, cmItem)
	return &types.Response{}, nil
}

func (l *NewChatMessageLogic) SendWSMessage(plist *[]model.Conversation, msgClientId string, cmItem model.ConversationMessage) {

	pc := types.WSPushMessage{
		WSPushBase: types.WSPushBase{
			Type:    defines.WSType_Message,
			SubType: defines.WSSubType_Message,
		},
		ChatId:          cmItem.ChatId,
		MessageId:       cmItem.Id,
		MessageClientId: msgClientId,
		MessageType:     cmItem.Type,
		MessageContent:  cmItem.Content,
		MessageSenderId: cmItem.SenderId,
		MessageSendTime: cmItem.SendTime,
	}

	pushString, _ := json.Marshal(pc)
	for _, conItem := range *plist {
		l.svcCtx.ChannelChat <- &types.ChannelMessage{
			UserId:  conItem.OwnerId,
			JwtId:   "",
			Content: pushString,
		}
	}
}

func (l *NewChatMessageLogic) SendPushNotificationToClient(plist *[]model.Conversation, cmItem model.ConversationMessage) {
	for _, person := range *plist {
		if person.OwnerId == cmItem.SenderId {
			continue
		}
		user, err := l.svcCtx.UserModel.FindOne(context.Background(), person.OwnerId)
		//fmt.Println("user: ", user, userId, err)
		if err != nil {
			return
		}

		devices, err := l.svcCtx.UserLoginRecordModel.QueryByUser(person.OwnerId)
		//fmt.Println("devices: ", devices)
		if err != nil {
			return
		}

		for _, item := range devices {
			if item.PushToken != "" {
				switch cmItem.Type {
				case defines.MsgType_Text:
					l.doPush(user.Name, cmItem.Content, item.PushToken)
				}
			}
		}
	}
}

func (l *NewChatMessageLogic) doPush(title string, content string, pushToken string) {
	logx.Info("SendPushToClient")

	certificate, _ := tls.LoadX509KeyPair(l.svcCtx.Config.Push.CERT, l.svcCtx.Config.Push.KEY)
	//l.svcCtx.Config.Push.CERT
	c, err := apns.NewClient(
		apns.WithCertificate(certificate),
		apns.WithBundleID("com.cabital.cabital.debug.h5.container"),
		apns.WithMaxIdleConnections(10),
		apns.WithTimeout(10*time.Second),
		apns.WithEndpoint("https://api.sandbox.push.apple.com:443"),
	)
	if err != nil {
		fmt.Println(err)
		return
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
