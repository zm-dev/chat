package db_store

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/zm-dev/chat/enum"
	"github.com/zm-dev/chat/model"
)

type dbCertificate struct {
	db *gorm.DB
}

const (
	_userCountGroupByTypeSQL = "SELECT `type` AS certificate_type, count(*) AS total FROM `certificates` GROUP BY `type`"
)

func (c *dbCertificate) CertificateCountGroupByType() ([]*model.CertificateCountResult, error) {
	rows, err := c.db.Raw(_userCountGroupByTypeSQL).Rows()
	if err != nil {
		err = errors.WithStack(err)
		return nil, err
	}
	defer rows.Close()

	res := make([]*model.CertificateCountResult, 0, 3)
	for rows.Next() {
		r := model.CertificateCountResult{}
		if err = c.db.ScanRows(rows, &r); err != nil {
			err = errors.WithStack(err)
			return nil, err
		}
		res = append(res, &r)
	}
	return res, err
}

func (c *dbCertificate) CertificateLoadByUserId(userId int64) (certificate *model.Certificate, err error) {
	if userId == 0 {
		return nil, model.ErrCertificateNotExist
	}
	certificate = &model.Certificate{}
	err = c.db.Where(model.Certificate{UserId: userId}).First(&certificate).Error
	if gorm.IsRecordNotFoundError(err) {
		err = model.ErrCertificateNotExist
	}
	return
}

func (c *dbCertificate) CertificateExist(account string) (bool, error) {
	var count uint8
	err := c.db.Model(&model.Certificate{}).Where(model.Certificate{Account: account}).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (c *dbCertificate) CertificateDelete(userId int64) error {
	return c.db.Where("user_id = ?", userId).Delete(&model.Certificate{}).Error
}

func (c *dbCertificate) CertificateIsNotExistErr(err error) bool {
	return model.CertificateIsNotExistErr(err)
}

func (c *dbCertificate) CertificateLoadByAccount(account string) (certificate *model.Certificate, err error) {
	if account == "" {
		return nil, model.ErrCertificateNotExist
	}
	certificate = &model.Certificate{}
	err = c.db.Where(model.Certificate{Account: account}).First(&certificate).Error
	if gorm.IsRecordNotFoundError(err) {
		err = model.ErrCertificateNotExist
	}
	return
}

func (c *dbCertificate) CertificateCreate(certificate *model.Certificate) error {
	return c.db.Create(certificate).Error
}

func (c *dbCertificate) CertificateUpdate(oldAccount, newAccount string, certificateType enum.CertificateType) error {
	return c.db.Model(&model.User{}).
		Where("account", oldAccount).
		Where("type", certificateType).
		UpdateColumn("account", newAccount).Error
}

func NewDBCertificate(db *gorm.DB) model.CertificateStore {
	return &dbCertificate{db: db}
}
