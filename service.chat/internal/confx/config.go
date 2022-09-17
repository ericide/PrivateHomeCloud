package confx

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"regexp"
	"strings"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc"
)

// DGTPrefix is a filter for environment variables
const ENVPrefix = "PHC_"

var loaders = map[string]func([]byte, interface{}) error{
	".yaml": conf.LoadConfigFromYamlBytes,
	".yml":  conf.LoadConfigFromYamlBytes,
	".json": conf.LoadConfigFromJsonBytes,
}

func LoadConfig(file string, v interface{}) error {
	kv := map[string]string{}

	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		if strings.HasPrefix(pair[0], ENVPrefix) {
			key := strings.TrimPrefix(pair[0], ENVPrefix)
			nkey := "@" + key + "@"
			// fmt.Printf("pair0[%s] nkey[%s]\n", pair[0], nkey)
			kv[nkey] = os.Getenv(pair[0])
		}
	}
	if content, err := ioutil.ReadFile(file); err != nil {
		return err
	} else if loader, ok := loaders[path.Ext(file)]; ok {
		return loader(globalReplace(content, kv), v)
	} else {
		return fmt.Errorf("unrecoginized file type: %s", file)
	}
}

func globalReplace(content []byte, kv map[string]string) []byte {
	res := string(content)
	for k, v := range kv {
		// fmt.Printf("key[%s] value[%s]\n", k, v)
		m := regexp.MustCompile(k)
		res = m.ReplaceAllString(res, v)
	}
	return []byte(res)
}

func MustLoad(path string, v interface{}) {
	if err := LoadConfig(path, v); err != nil {
		logx.Errorf("error: config file %s, %s", path, err.Error())
		os.Exit(1)
	}

	// reflect to discover version
	s := reflect.ValueOf(v).Elem()
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		switch f.Type() {
		case reflect.TypeOf(VersionConf{}):
			f.Field(0).SetString(LoadVersion())
		case reflect.TypeOf(zrpc.RpcServerConf{}):
			err := f.Interface().(zrpc.RpcServerConf).SetUp()
			if err != nil {
				logx.Error(err)
			}
		}
	}
}
