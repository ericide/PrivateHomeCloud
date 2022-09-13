package logic

import (
	"context"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"reflect"

	"service.chat/internal/svc"
	"service.chat/internal/types"
)

func JwtAuth(svc *svc.ServiceContext, tokenStr string) (*types.LoginClaims, error) {

	token, err := jwt.Parse(
		tokenStr,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(svc.Config.Auth.AccessSecret), nil
		},
	)

	if !token.Valid {
		return nil, errors.New("token invalid")
	}

	claim, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("claim type error")
	}

	if err != nil {
		return nil, err
	}

	fmt.Println(reflect.TypeOf(claim["exp"]))

	retClaim := types.LoginClaims{
		UserId: claim["user_id"].(string),
		JwtId:  claim["jwt_id"].(string),
	}

	return &retClaim, nil
}

func JwtClaimAuth(svc *svc.ServiceContext, claim *types.LoginClaims) error {

	_, err := svc.UserModel.FindOne(claim.UserId)
	if err != nil {
		return errors.New("user not exist")
	}

	record, err := svc.UserLoginRecordModel.FindOne(context.Background(), claim.JwtId)
	if err != nil {
		return errors.New("jwt is invalid")
	}

	fmt.Println(record.Invalid)

	if record.Invalid == 1 {
		return errors.New("jwt is invalid")
	}

	return nil
}
