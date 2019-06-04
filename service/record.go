package service

import (
	"context"
	"github.com/zm-dev/chat/model"
)

type recordService struct {
	rs model.RecordStore
}

func (rSvc *recordService) BatchSetRead(ids []int64, toId int64) error {
	return rSvc.rs.BatchSetRead(ids, toId)
}

func (rSvc *recordService) PageRecord(page *model.Page, teacherId, studentId int64) (err error) {
	return rSvc.rs.PageRecord(page, teacherId, studentId)
}

func (rSvc *recordService) CreateRecord(record *model.Record) error {
	return rSvc.rs.CreateRecord(record)
}

func BatchSetRead(ctx context.Context, ids []int64, toId int64) error {
	return FromContext(ctx).BatchSetRead(ids, toId)
}

func PageRecord(ctx context.Context, page *model.Page, teacherId, studentId int64) (err error) {
	return FromContext(ctx).PageRecord(page, teacherId, studentId)
}

func CreateRecord(ctx context.Context, record *model.Record) error {
	return FromContext(ctx).CreateRecord(record)
}

func NewRecordService(rs model.RecordStore) model.RecordService {
	return &recordService{rs: rs}
}
