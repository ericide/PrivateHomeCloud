package middleware

import (
	"fmt"
	"net/http"
	"service.chat/internal/model"
)

type UsercheckMiddleware struct {
	userModel            model.UserModel
	userLoginRecordModel model.UserLoginRecordModel
}

func NewUsercheckMiddleware(userModel model.UserModel, userLoginRecordModel model.UserLoginRecordModel) *UsercheckMiddleware {
	return &UsercheckMiddleware{
		userModel:            userModel,
		userLoginRecordModel: userLoginRecordModel,
	}
}

func (m *UsercheckMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		
		userId, ok := r.Context().Value("user_id").(string)
		if !ok {
			w.WriteHeader(401)
			w.Write([]byte(`{"message":"need login"}`))
			return
		}

		jwtId, ok := r.Context().Value("jwt_id").(string)
		if !ok {
			w.WriteHeader(401)
			w.Write([]byte(`{"message":"need login"}`))
			return
		}

		_, err := m.userModel.FindOne(userId)
		if err != nil {
			w.WriteHeader(401)
			w.Write([]byte(`{"message":"user not exist"}`))
			return
		}

		record, err := m.userLoginRecordModel.FindOne(r.Context(), jwtId)
		if err != nil {
			w.WriteHeader(401)
			w.Write([]byte(`{"message":"jwt is invalid"}`))
			return
		}

		fmt.Println(record.Invalid)

		if record.Invalid == 1 {
			w.WriteHeader(401)
			w.Write([]byte(`{"message":"jwt is invalid"}`))
			return
		}

		// Passthrough to next handler if need
		next(w, r)
	}
}
