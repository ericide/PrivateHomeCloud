package types

import (
	"time"
)

type LoginClaims struct {
	UserId string `json:"user_id"`
	JwtId  string `json:"jwt_id"`
	Exp    int64  `json:"exp"`
	Iat    int64  `json:"iat"`
}
type ChannelMessage struct {
	UserId  string
	JwtId   string
	Content []byte
}

type WSPushMessage struct {
	Type              string    `json:"user_id"`
	MessageChatId     string    `json:"message_chat_id"`
	MessageType       string    `json:"message_type"`
	MessageContent    string    `json:"message_content"`
	MessageSenderId   string    `json:"message_sender_id"`
	MessageCreateTime time.Time `json:"message_create_time"`
}
