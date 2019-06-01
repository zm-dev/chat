package model

import (
	"errors"
	"time"
)

type User struct {
	Id         int64     `gorm:"type:BIGINT AUTO_INCREMENT;PRIMARY_KEY;NOT NUll" json:"id"`
	Email      string    `gorm:"type:varchar(40)" json:"email"`             // 邮箱
	AvatarHash string    `gorm:"type:char(32)" json:"avatar_hash"`          // 头像
	NikeName   string    `gorm:"type:varchar(50)" json:"nike_name"`         // 昵称
	Profile    string    `gorm:"type:varchar(255)" json:"profile"`          // 简介
	Gender     uint8     `gorm:"type:TINYINT;DEFAULT:0" json:"gender"`      // 性别
	Password   string    `gorm:"type:varchar(64);NOT NULL" json:"password"` // 密码
	Status     uint8     `gorm:"type:TINYINT;DEFAULT:0" json:"status"`      // 在线状态
	PwPlain    string    `gorm:"type:varchar(20);not null" json:"-"`        // 密码明文
	IsAdmin    bool      `gorm:"type:TINYINT" json:"is_admin"`              // 管理员
	IsTeacher  bool      `gorm:"type:TINYINT" json:"is_teacher"`            // 教师
	IsStudent  bool      `gorm:"type:TINYINT" json:"is_student"`            // 学生
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type UserType int8

const (
	TeacherType UserType = iota
	StudentType
	AdminType
)

const (
	GenderMan = iota
	GenderWoman
	GenderSecrecy
)

const (
	UserStatusOnline = iota
	UserStatusOffline
)

var ErrUserNotExist = errors.New("user is not exist")

var ErrUserTypeNotExist = errors.New("user type has not allow")

type UserStore interface {
	UserLoad(int64) (*User, error)
	UserUpdate(*User) error
	UserCreate(*User) error
	UserList(uType UserType) ([]*User, error)
}

type UserService interface {
	UserStore
	UserLogin(account, password string) (*Ticket, error)
	UserRegister(account string, certificateType CertificateType, password string) (userId int64, err error)
	UserUpdatePassword(userId int64, newPassword string) (err error)
	TeacherList() ([]*User, error)
	StudentList() ([]*User, error)
}
