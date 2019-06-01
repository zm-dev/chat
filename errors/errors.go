package errors

import (
	"errors"
	"github.com/zm-dev/gerrors"
)

var RecordNotFoundError = errors.New("record not found")

// 参数绑定出错
func BindError(err error) error {
	// todo 处理422
	return gerrors.BadRequest(10001, err.Error(), err)
}

func BadRequest(msg string, err ...error) error {
	return gerrors.BadRequest(10002, msg, err...)
}

func InternalServerError(msg string, err ...error) error {
	return gerrors.InternalServerError(10003, msg, err...)
}

func Unauthorized(message ...string) error {
	var msg string
	if len(message) == 0 {
		msg = "请先登录"
	} else {
		msg = message[0]
	}
	return gerrors.Unauthorized(10004, msg, nil)
}

// NotFound generates a 404 error.
func NotFound(message string, err ...error) error {
	return gerrors.NotFound(10005, message, err...)
}

// 没有权限
func Forbidden(msg string, err ...error) error {
	return gerrors.Forbidden(10006, msg, err...)
}

func ErrAccountAlreadyExisted() error {
	return gerrors.BadRequest(10007, "account already existed", nil)
}

func ErrPassword() error {
	return gerrors.BadRequest(10008, "密码错误", nil)
}

func ErrAccountNotFound() error {
	return gerrors.NotFound(10009, "账号不存在", nil)
}
