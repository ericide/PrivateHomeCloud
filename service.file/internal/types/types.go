package types

import "errors"

type DataResponse struct {
	Data any `json:"data"`
}

type RespFile struct {
	Name    string `json:"name"`
	Size    int64  `json:"size"`
	ModTime string `json:"mod_time"`
	IsDir   bool   `json:"is_dir"`
}

var (
	// ErrNotFound is an alias of sqlx.ErrNotFound.
	ErrNotFound = errors.New("Not FOUND")
)