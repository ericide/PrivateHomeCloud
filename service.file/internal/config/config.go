package config

import (
	"github.com/golang-jwt/jwt"
	"golang.org/x/net/webdav"
	"net/http"
)

type Config struct {
	Port         int
	AccessToken  string
	PhysicalPath string
	Webdav       *webdav.Handler
}

func (c *Config) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	jwtString := r.Header.Get("authorization")

	_, err := jwt.Parse(jwtString, func(token *jwt.Token) (interface{}, error) {
		return []byte(c.AccessToken), nil
	})
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	//fmt.Println(token, r.RequestURI, r.Method, r.Header.Get("Depth"))

	c.Webdav.ServeHTTP(w, r)

}
