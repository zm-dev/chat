package model

import (
	"errors"
	"github.com/zm-dev/chat/enum"
	"time"
)

type User struct {
	Id         int64       `gorm:"type:BIGINT AUTO_INCREMENT;PRIMARY_KEY;NOT NUll" json:"id"`
	AvatarHash string      `gorm:"type:char(32)" json:"avatar_hash"`          // 头像
	NickName   string      `gorm:"type:varchar(50)" json:"nick_name"`         // 昵称
	Profile    string      `gorm:"type:varchar(255)" json:"profile"`          // 简介
	Gender     enum.Gender `gorm:"type:TINYINT;DEFAULT:0" json:"gender"`      // 性别
	GroupId    enum.Group  `gorm:"type:TINYINT;DEFAULT:0" json:"group"`       // 组
	Company    string      `gorm:"type:varchar(50)" json:"company"`           // 工作单位
	Password   string      `gorm:"type:varchar(64);NOT NULL" json:"password"` // 密码
	PwPlain    string      `gorm:"type:varchar(20);not null" json:"-"`        // 密码明文
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
}

var ErrUserNotExist = errors.New("user is not exist")

var ErrUserTypeNotExist = errors.New("user type has not allow")

type UserStore interface {
	UserLoad(int64) (*User, error)
	UserUpdate(*User) error
	UserCreate(*User) error
	UserList(uType enum.CertificateType, page *Page) error
}

type UserService interface {
	UserStore
	UserLogin(account, password string) (*Ticket, error)
	UserRegister(account, password, nickname string, certificateType enum.CertificateType) (userId int64, err error)
	UserUpdatePassword(userId int64, newPassword string) (err error)
	TeacherList() ([]*User, error)
	StudentList(*Page) error
}
