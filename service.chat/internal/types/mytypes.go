package types

var FormatISOTime = "2006-01-02T15:04:05.000Z07:00"

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

type WSPushBase struct {
	Type    string `json:"type"`
	SubType string `json:"sub_type"`
}

type WSPushMessage struct {
	WSPushBase
	ChatId            string `json:"chat_id"`
	MessageType       string `json:"message_type"`
	MessageId         string `json:"message_id"`
	MessageClientId   string `json:"message_client_id"`
	MessageContent    string `json:"message_content"`
	MessageSenderId   string `json:"message_sender_id"`
	MessageCreateTime string `json:"message_create_time"`
}
type WSUpdateReadTime struct {
	WSPushBase
	ChatId   string `json:"chat_id"`
	UserId   string `json:"user_id"`
	ReadTime string `json:"read_time"`
}
