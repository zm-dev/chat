package service

import (
	"github.com/zm-dev/chat_v2/model"
)

type certificateService struct {
	model.CertificateStore
}

func NewCertificateService(cs model.CertificateStore) model.CertificateService {
	return &certificateService{cs}
}
