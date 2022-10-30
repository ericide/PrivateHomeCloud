package config

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type Config struct {
	Port        int
	AccessToken string
}

func (c *Config) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	fmt.Println(r.RequestURI)
	req, _ := http.NewRequest(r.Method, "http://localhost:8001"+r.RequestURI, r.Body)
	for k, v := range r.Header {
		for _, vv := range v {
			req.Header.Add(k, vv)
		}
	}
	req.Header.Add("authorization", "bbbbb")

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		print(err.Error())
		return
	}

	for k, v := range resp.Header {
		for _, vv := range v {
			w.Header().Add(k, vv)
		}
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error:", err)
	}

	_, err = w.Write(result)
	if err != nil {
		fmt.Println("error:", err)
	}
	//jwtString := r.Header.Get("authorization")
	//
	//token, err := jwt.Parse(jwtString, func(token *jwt.Token) (interface{}, error) {
	//	return []byte(c.AccessToken), nil
	//})
	//if err != nil {
	//	w.WriteHeader(http.StatusUnauthorized)
	//	return
	//}

	//fmt.Println(token, r.RequestURI, r.Method, r.Header.Get("Depth"))

}
