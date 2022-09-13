package logic

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"service.chat/internal/model"
	"time"

	"service.chat/internal/svc"
	"service.chat/internal/types"

	"github.com/edganiukov/apns"
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
	userId, ok := l.ctx.Value("user_id").(string)
	if !ok {
		return nil, errors.New("user id not exist")
	}

	list, err := l.svcCtx.ConversationModel.QueryByChatId(req.ChatId)

	fmt.Println("list: ", list)

	if err != nil {
		return nil, err
	}
	if len(*list) == 0 {
		fmt.Println("errors: ", list)
		return nil, errors.New("chat id is invalid")
	}

	cmItem := model.ConversationMessage{
		Id:         uuid.New().String(),
		ChatId:     req.ChatId,
		Type:       req.Type,
		SenderId:   userId,
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
	for _, conItem := range *list {
		l.SendMessage(conItem.OwnerId, "", pushString)
	}

	go l.SendPushToClient(cmItem.Content, "cd226ff5c9b9b72dc7347c7ea81bde081a95e91913f382dff4b1432a24fd47e0")

	return &types.Response{}, nil
}

func (l *NewChatMessageLogic) SendMessage(userId string, jwtId string, content []byte) {
	l.svcCtx.ChannelChat <- &types.ChannelMessage{
		UserId:  userId,
		JwtId:   jwtId,
		Content: content,
	}
}

func (l *NewChatMessageLogic) SendPushToClient(pushText string, pushToken string) {
	fmt.Println("SendPushToClient")

	data, err := ioutil.ReadFile(l.svcCtx.Config.Push.KEY)
	if err != nil {
		log.Fatal(err)
	}
	//l.svcCtx.Config.Push.CERT
	c, err := apns.NewClient(
		apns.WithCertificate(tls.Certificate{}),
		apns.WithBundleID("com.cabital.cabital.debug.h5.container"),
		apns.WithMaxIdleConnections(10),
		apns.WithTimeout(5*time.Second),
	)
	if err != nil {
		/* ... */
	}
	resp, err := c.Send("<device token>",
		apns.Payload{
			APS: apns.APS{
				Alert: apns.Alert{
					Title: "Test Push",
					Body:  "Hi world",
				},
			},
		},
		apns.WithExpiration(10),
		apns.WithCollapseID("test-collapse-id"),
		apns.WithPriority(5),
	)

	if err != nil {
		/* ... */
	}

}
