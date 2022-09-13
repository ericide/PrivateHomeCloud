package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ UserLoginRecordModel = (*customUserLoginRecordModel)(nil)

type (
	// UserLoginRecordModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserLoginRecordModel.
	UserLoginRecordModel interface {
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
