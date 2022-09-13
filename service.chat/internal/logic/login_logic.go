package logic

import (
	"context"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"service.chat/internal/model"
	"time"

	"service.chat/internal/svc"
	"service.chat/internal/types"

	"github.com/dgrijalva/jwt-go"
	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) LoginLogic {
	return LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req types.LoginReq) (*types.LoginResponse, error) {
	// todo: add your logic here and delete this line

	userInfo, err := l.svcCtx.UserModel.FindOneByPhone(l.ctx, req.Username)

	if err != nil {
		return nil, err
	}

	if userInfo.Password != req.Password {
		return nil, errors.New("invalid username or password")
	}

	//---start---
	now := time.Now().Unix()
	accessExpire := l.svcCtx.Config.Auth.AccessExpire
	jwtToken, err := l.getJwtToken(l.svcCtx.Config.Auth.AccessSecret, now, l.svcCtx.Config.Auth.AccessExpire, userInfo.Id)
	if err != nil {
		return nil, err
	}

	return &types.LoginResponse{
		UserId:       userInfo.Id,
		AccessToken:  jwtToken,
		AccessExpire: now + accessExpire,
		RefreshAfter: now + accessExpire/2,
	}, nil
}

func (l *LoginLogic) getJwtToken(secretKey string, iat, seconds int64, userId string) (string, error) {
	jwtId := uuid.New().String()
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds
	claims["iat"] = iat
	claims["user_id"] = userId
	claims["jwt_id"] = jwtId
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims

	_, err := l.svcCtx.UserLoginRecordModel.Insert(l.ctx, &model.UserLoginRecord{
		Id:         jwtId,
		UserId:     userId,
		Device:     "random",
		DeviceName: "random random",
		PushToken:  "",
		Invalid:    0,
		CreateTime: time.Now(),
	})

	if err != nil {
		return "", err
	}

	return token.SignedString([]byte(secretKey))
}
