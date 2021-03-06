package service

import (
	"context"

	"github.com/zm-dev/chat/model"
)

type recordService struct {
	rs model.RecordStore
}

func (rSvc *recordService) GetNotReadRecordCount(fromId, toId int64) int32 {
	return rSvc.rs.GetNotReadRecordCount(fromId, toId)
}

func (rSvc *recordService) LastRecordList(userIdA int64, size int) (records []*model.Record, err error) {
	return rSvc.rs.LastRecordList(userIdA, size)
}

func (rSvc *recordService) BatchSetRead(ids []int64, toId int64) error {
	return rSvc.rs.BatchSetRead(ids, toId)
}

func (rSvc *recordService) SetAllRead(fromId, toId int64) error {
	return rSvc.rs.SetAllRead(fromId, toId)
}

func (rSvc *recordService) PageRecord(page *model.Page, userIdA, userIdB int64, onlyShowNotRead, isOrderAsc bool) (err error) {
	return rSvc.rs.PageRecord(page, userIdA, userIdB, onlyShowNotRead, isOrderAsc)
}

func (rSvc *recordService) CreateRecord(record *model.Record) (int64, error) {
	return rSvc.rs.CreateRecord(record)
}

func SetAllRead(ctx context.Context, fromId, toId int64) error {
	return FromContext(ctx).SetAllRead(fromId, toId)
}

func BatchSetRead(ctx context.Context, ids []int64, toId int64) error {
	return FromContext(ctx).BatchSetRead(ids, toId)
}

func PageRecord(ctx context.Context, page *model.Page, userIdA, userIdB int64, onlyShowNotRead, isOrderAsc bool) (err error) {
	return FromContext(ctx).PageRecord(page, userIdA, userIdB, onlyShowNotRead, isOrderAsc)
}

func CreateRecord(ctx context.Context, record *model.Record) (int64, error) {
	return FromContext(ctx).CreateRecord(record)
}

func LastRecordList(ctx context.Context, userIdA int64, size int) (records []*model.Record, err error) {
	return FromContext(ctx).LastRecordList(userIdA, size)
}

func GetNotReadRecordCount(ctx context.Context, fromId, toId int64) int32 {
	return FromContext(ctx).GetNotReadRecordCount(fromId, toId)
}

func NewRecordService(rs model.RecordStore) model.RecordService {
	return &recordService{rs: rs}
}
