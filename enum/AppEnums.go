package enum

type UserType uint8

type Gender uint8

type Group uint8

type CertificateType uint8

const (
	CertificateAdmin   CertificateType = iota // 管理员
	CertificateStudent                        // 学生
	CertificateTeacher                        // 老师
)

// 性别枚举
const (
	GenderMan     Gender = iota // 男
	GenderWoman                 // 女
	GenderSecrecy               // 保密
)

// 群组枚举
const (
	TeacherGroup Group = iota // 老师
	AlumnusGroup              // 校友
	PBFGroup                  // 朋辈辅导员
)

func ParseGroup(g Group) (result string) {
	switch g {
	case TeacherGroup:
		result = "老师"
	case AlumnusGroup:
		result = "校友"
	case PBFGroup:
		result = "朋辈辅导员"
	default:
		result = ""
	}
	return result
}

func ParseGender(g Gender) (gender string) {
	switch g {
	case GenderMan:
		gender = "男"
	case GenderWoman:
		gender = "女"
	case GenderSecrecy:
		fallthrough
	default:
		gender = "未知"
	}
	return
}
