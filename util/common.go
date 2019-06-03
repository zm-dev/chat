package util

import "github.com/zm-dev/chat/model"

// 转换用户性别
func ConvertUserGender(g uint8) (gender string) {
	switch g {
	case model.GenderMan:
		gender = "男"
	case model.GenderWoman:
		gender = "女"
	case model.GenderSecrecy:
		fallthrough
	default:
		gender = "未知"
	}
	return
}

// 转换用户在线状态
//func ConvertUserStatus(s uint8) (status string) {
//	switch s {
//	case model.UserStatusOnline:
//		status = "在线"
//	case model.UserStatusOffline:
//		status = "离线"
//	default:
//		status = ""
//	}
//	return
//}
