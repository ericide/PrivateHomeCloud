package confx

import (
	"io/ioutil"
	"log"
	"strings"
)

const (
	defaultVersion = "NULL"
	versionFile    = "version.txt"
)

func LoadVersion() string {
	content, err := ioutil.ReadFile(versionFile)
	if err != nil {
		log.Println("warning -- no version avialable")
		log.Print(err)
	}
	if len(content) > 0 {
		return strings.TrimSuffix(string(content), "\n")
	} else {
		return defaultVersion
	}
}
