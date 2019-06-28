package model

import (
	"errors"
	"github.com/zm-dev/chat/enum"
)

type Certificate struct {
	Id      int64                `gorm:"type:BIGINT AUTO_INCREMENT;PRIMARY_KEY;NOT NULL"`
	UserId  int64                `gorm:"type:BIGINT;INDEX"` // 账号详细信息
	Account string               `gorm:"NOT NULL;UNIQUE"`   // 账户名称（教师：工号；学生：学号；管理员：username）
	Type    enum.CertificateType `gorm:"type:TINYINT"`      // 账号类型（教师；学生；管理员）
}

type CertificateCountResult struct {
	CertificateType enum.CertificateType
	Total           int64
}

type CertificateStore interface {
	CertificateExist(account string) (bool, error)
	CertificateDelete(userId int64) error
	CertificateLoadByAccount(account string) (*Certificate, error)
	CertificateLoadByUserId(userId int64) (*Certificate, error)
	CertificateIsNotExistErr(error) bool
	CertificateCreate(certificate *Certificate) error
	CertificateUpdate(oldAccount, newAccount string, certificateType enum.CertificateType) error
	CertificateCountGroupByType() ([]*CertificateCountResult, error)
}

var ErrCertificateNotExist = errors.New("certificate not exist")

func CertificateIsNotExistErr(err error) bool {
	return err == ErrCertificateNotExist
}

type CertificateService interface {
	CertificateStore
	UserCountWithCertificate() (map[string]int64, error)
}
