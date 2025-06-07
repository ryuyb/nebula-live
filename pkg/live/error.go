package live

import "errors"

var (
	ErrRoomNotExist  = errors.New("room not exists")
	ErrInternalError = errors.New("internal error")
)
