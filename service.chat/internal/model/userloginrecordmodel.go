package model

import (
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UserLoginRecordModel = (*customUserLoginRecordModel)(nil)

type (
	// UserLoginRecordModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserLoginRecordModel.
	UserLoginRecordModel interface {
		QueryByUser(userId string) ([]*UserLoginRecord, error)
		userLoginRecordModel
	}

	customUserLoginRecordModel struct {
		*defaultUserLoginRecordModel
	}
)

// NewUserLoginRecordModel returns a model for the database table.
func NewUserLoginRecordModel(conn sqlx.SqlConn) UserLoginRecordModel {
	return &customUserLoginRecordModel{
		defaultUserLoginRecordModel: newUserLoginRecordModel(conn),
	}
}

func (m *customUserLoginRecordModel) QueryByUser(userId string) ([]*UserLoginRecord, error) {
	query := fmt.Sprintf("select %s from %s where `user_id` = ? ", userLoginRecordRows, m.table)
	var resp []*UserLoginRecord
	err := m.conn.QueryRows(&resp, query, userId)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return []*UserLoginRecord{}, nil
	default:
		return nil, err
	}
}
