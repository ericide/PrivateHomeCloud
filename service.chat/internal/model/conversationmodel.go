package model

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
	"github.com/zeromicro/go-zero/tools/goctl/model/sql/builderx"
)

var (
	conversationFieldNames          = builderx.RawFieldNames(&Conversation{})
	conversationRows                = strings.Join(conversationFieldNames, ",")
	conversationRowsExpectAutoSet   = strings.Join(stringx.Remove(conversationFieldNames, "`create_time`", "`update_time`"), ",")
	conversationRowsWithPlaceHolder = strings.Join(stringx.Remove(conversationFieldNames, "`id`", "`create_time`", "`update_time`"), "=?,") + "=?"
)

type (
	ConversationModel interface {
		Insert(data Conversation) (sql.Result, error)
		FindOne(id string) (*Conversation, error)
		QueryByChatId(id string) (*[]Conversation, error)
		Update(data Conversation) error
		Delete(id string) error
		Range(userId string) ([]*Conversation, error)
	}

	defaultConversationModel struct {
		conn  sqlx.SqlConn
		table string
	}

	Conversation struct {
		Id           string    `db:"id"`
		Type         string    `db:"type"`
		ChatId       string    `db:"chat_id"`
		OwnerId      string    `db:"owner_id"`
		Name         string    `db:"name"`
		LastReadTime time.Time `db:"last_read_time"`
		CreateTime   time.Time `db:"create_time"`
	}
)

func NewConversationModel(conn sqlx.SqlConn) ConversationModel {
	return &defaultConversationModel{
		conn:  conn,
		table: "`conversation`",
	}
}

func (m *defaultConversationModel) Insert(data Conversation) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?)", m.table, conversationRowsExpectAutoSet)
	ret, err := m.conn.Exec(query, data.Id, data.Type, data.ChatId, data.OwnerId, data.Name, data.LastReadTime)
	return ret, err
}

func (m *defaultConversationModel) FindOne(id string) (*Conversation, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", conversationRows, m.table)
	var resp Conversation
	err := m.conn.QueryRow(&resp, query, id)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
func (m *defaultConversationModel) QueryByChatId(id string) (*[]Conversation, error) {
	query := fmt.Sprintf("select %s from %s where `chat_id` = ?", conversationRows, m.table)
	var resp []Conversation
	err := m.conn.QueryRows(&resp, query, id)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return &[]Conversation{}, nil
	default:
		return nil, err
	}
}

func (m *defaultConversationModel) Update(data Conversation) error {
	query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, conversationRowsWithPlaceHolder)
	_, err := m.conn.Exec(query, data.Type, data.ChatId, data.OwnerId, data.Name, data.LastReadTime, data.Id)
	return err
}

func (m *defaultConversationModel) Delete(id string) error {
	query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
	_, err := m.conn.Exec(query, id)
	return err
}

func (m *defaultConversationModel) Range(userId string) ([]*Conversation, error) {
	query := fmt.Sprintf("select %s from %s where `owner_id` = ? ", conversationRows, m.table)
	var resp []*Conversation
	err := m.conn.QueryRows(&resp, query, userId)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return []*Conversation{}, nil
	default:
		return nil, err
	}
}
