package model

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ConversationModel = (*customConversationModel)(nil)

type (
	// ConversationModel is an interface to be customized, add more methods here,
	// and implement the added methods in customConversationModel.
	ConversationModel interface {
		QueryByChatId(ctx context.Context, chatId string) (*[]Conversation, error)
		QueryByUserId(ctx context.Context, userId string) (*[]Conversation, error)
		QueryByChatIdUserId(ctx context.Context, chatId string, userId string) (*Conversation, error)
		conversationModel
	}

	customConversationModel struct {
		*defaultConversationModel
	}
)

// NewConversationModel returns a model for the database table.
func NewConversationModel(conn sqlx.SqlConn) ConversationModel {
	return &customConversationModel{
		defaultConversationModel: newConversationModel(conn),
	}
}
func (m *customConversationModel) QueryByChatId(ctx context.Context, chatId string) (*[]Conversation, error) {
	query := fmt.Sprintf("select %s from %s where `chat_id` = ?", conversationRows, m.table)
	var resp []Conversation
	err := m.conn.QueryRowsCtx(ctx, &resp, query, chatId)
	fmt.Println(err)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return &[]Conversation{}, nil
	default:
		return nil, err
	}
}

func (m *customConversationModel) QueryByChatIdUserId(ctx context.Context, chatId string, userId string) (*Conversation, error) {
	query := fmt.Sprintf("select %s from %s where `chat_id` = ? AND `owner_id` = ?", conversationRows, m.table)
	var resp Conversation
	err := m.conn.QueryRowCtx(ctx, &resp, query, chatId, userId)
	fmt.Println(err)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *customConversationModel) QueryByUserId(ctx context.Context, userId string) (*[]Conversation, error) {
	query := fmt.Sprintf("select %s from %s where `owner_id` = ?", conversationRows, m.table)
	var resp []Conversation
	err := m.conn.QueryRowsCtx(ctx, &resp, query, userId)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return &[]Conversation{}, nil
	default:
		return nil, err
	}
}
