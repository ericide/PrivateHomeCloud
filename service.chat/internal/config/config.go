package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf
	Mysql struct {
		DataSource string
	}
	Auth struct {
		AccessSecret string
		AccessExpire int64
	}
	Push struct {
		IOS struct {
			CERT     string
			KEY      string
			BundleId string
		}
	}
}
