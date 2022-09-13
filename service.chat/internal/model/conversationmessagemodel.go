package model

import (
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ConversationMessageModel = (*customConversationMessageModel)(nil)

type (
	// ConversationMessageModel is an interface to be customized, add more methods here,
	// and implement the added methods in customConversationMessageModel.
	ConversationMessageModel interface {
		Range(chatId string, start int, length int) ([]*ConversationMessage, error)
		conversationMessageModel
	}

	customConversationMessageModel struct {
		*defaultConversationMessageModel
	}
)

// NewConversationMessageModel returns a model for the database table.
func NewConversationMessageModel(conn sqlx.SqlConn) ConversationMessageModel {
	return &customConversationMessageModel{
		defaultConversationMessageModel: newConversationMessageModel(conn),
	}
}

func (m *customConversationMessageModel) Range(chatId string, start int, length int) ([]*ConversationMessage, error) {
	query := fmt.Sprintf("select %s from %s where `chat_id` = ? order by `create_time` desc limit ?, ? ", conversationMessageRows, m.table)
	var resp []*ConversationMessage
	err := m.conn.QueryRows(&resp, query, chatId, start, length)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return []*ConversationMessage{}, nil
	default:
		return nil, err
	}
}
