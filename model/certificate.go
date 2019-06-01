package model

import (
	"errors"
)

type CertificateType uint8

const (
	CertificateUserName CertificateType = iota
	CertificatePhoneNum
	CertificateEmail
)

type Certificate struct {
	Id      int64           `gorm:"type:BIGINT AUTO_INCREMENT;PRIMARY_KEY;NOT NULL"`
	UserId  int64           `gorm:"type:BIGINT;INDEX"` // 账号详细信息
	Account string          `gorm:"NOT NULL;UNIQUE"`   // 账户名称（教师：工号；学生：学号；管理员：username）
	Type    CertificateType `gorm:"type:TINYINT"`      // 账号类型
}

type CertificateStore interface {
	CertificateExist(account string) (bool, error)
	CertificateLoadByAccount(account string) (*Certificate, error)
	CertificateIsNotExistErr(error) bool
	CertificateCreate(certificate *Certificate) error
	CertificateUpdate(oldAccount, newAccount string, certificateType CertificateType) error
}

var ErrCertificateNotExist = errors.New("certificate not exist")

func CertificateIsNotExistErr(err error) bool {
	return err == ErrCertificateNotExist
}

type CertificateService interface {
	CertificateStore
}
