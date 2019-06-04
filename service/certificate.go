package service

import (
	"context"
	"github.com/zm-dev/chat/model"
)

type certificateService struct {
	model.CertificateStore
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

func NewCertificateService(cs model.CertificateStore) model.CertificateService {
	return &certificateService{cs}
}
