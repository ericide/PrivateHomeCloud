package types

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
	Message RespConversationMessage `json:"message"`
}
type WSUpdateReadTime struct {
	WSPushBase
	ChatId   string `json:"chat_id"`
	UserId   string `json:"user_id"`
	ReadTime int64  `json:"read_time"`
}

type WSBaseReq struct {
	Type string `json:"type"`
}

type WSAuthReq struct {
	Token string `json:"token"`
}
