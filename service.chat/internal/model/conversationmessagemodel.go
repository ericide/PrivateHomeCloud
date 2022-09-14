package model

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"time"
)

var _ ConversationMessageModel = (*customConversationMessageModel)(nil)

type (
	// ConversationMessageModel is an interface to be customized, add more methods here,
	// and implement the added methods in customConversationMessageModel.
	ConversationMessageModel interface {
		Range(chatId string, start int, length int) ([]*ConversationMessage, error)
		CountAfterTime(ctx context.Context, chatId string, time time.Time) (*int, error)
		LastMessage(ctx context.Context, chatId string) (*ConversationMessage, error)
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

func (m *customConversationMessageModel) CountAfterTime(ctx context.Context, chatId string, time time.Time) (*int, error) {
	query := fmt.Sprintf("select Count(*) from %s where `chat_id` = ? AND `create_time` > ?", m.table)
	var count int
	err := m.conn.QueryRowCtx(ctx, &count, query, chatId, time)
	switch err {
	case nil:
		return &count, nil
	case sqlc.ErrNotFound:
		count = 0
		return &count, nil
	default:
		return nil, err
	}
}
func (m *customConversationMessageModel) LastMessage(ctx context.Context, chatId string) (*ConversationMessage, error) {
	query := fmt.Sprintf("select %s from %s where `chat_id` = ? order by create_time desc limit 1", conversationMessageRows, m.table)
	var resp ConversationMessage
	err := m.conn.QueryRowCtx(ctx, &resp, query, chatId)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
