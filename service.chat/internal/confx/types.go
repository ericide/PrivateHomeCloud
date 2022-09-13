package confx

import "os"

type EnvType int

const (
	EnvTypeLocal EnvType = iota + 1
	EnvTypeQA
	EnvTypeLive
)

const EnvKey = "ENV"

func (w EnvType) String() string {
	switch w {
	case EnvTypeLocal:
		return "local"
	case EnvTypeQA:
		return "qa"
	default:
		return "live"
	}
}

func MatchEnv(types ...EnvType) bool {
	if value, ok := os.LookupEnv(EnvKey); ok {
		for _, t := range types {
			if t.String() == value {
				return true
			}
		}
	}
	return false
}

type VersionConf struct {
	Version string `json:",optional"` // Version of services
}
