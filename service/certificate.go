package service

import (
	"context"
	"github.com/zm-dev/chat/enum"
	"github.com/zm-dev/chat/model"
)

type certificateService struct {
	model.CertificateStore
}

func (cSvc *certificateService) UserCountWithCertificate() (map[string]int64, error) {
	res, err := cSvc.CertificateStore.CertificateCountGroupByType()
	if err != nil {
		return nil, err
	}
	userGroupCount := make(map[string]int64, 3)
	for _, v := range res {
		switch v.CertificateType {
		case enum.CertificateAdmin:
			userGroupCount["admin_count"] = v.Total
		case enum.CertificateStudent:
			userGroupCount["student_count"] = v.Total
		case enum.CertificateTeacher:
			userGroupCount["teacher_count"] = v.Total
		default:
			continue
		}
	}

	return userGroupCount, nil
}

func (cSvc *certificateService) CertificateLoadByAccount(account string) (*model.Certificate, error) {
	return cSvc.CertificateStore.CertificateLoadByAccount(account)
}

func CertificateLoadByAccount(c context.Context, account string) (*model.Certificate, error) {
	return FromContext(c).CertificateLoadByAccount(account)
}

func (cSvc *certificateService) CertificateLoadByUserId(userId int64) (*model.Certificate, error) {
	return cSvc.CertificateStore.CertificateLoadByUserId(userId)
}

func CertificateLoadByUserId(c context.Context, userId int64) (*model.Certificate, error) {
	return FromContext(c).CertificateLoadByUserId(userId)
}

func UserCountWithCertificate(c context.Context) (map[string]int64, error) {
	return FromContext(c).UserCountWithCertificate()
}

func NewCertificateService(cs model.CertificateStore) model.CertificateService {
	return &certificateService{cs}
}
