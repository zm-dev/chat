package db_store

import (
	"github.com/jinzhu/gorm"
	"github.com/zm-dev/chat/model"
	"time"
)

type dbRecord struct {
	db *gorm.DB
}

func (r *dbRecord) LastRecordList(toId int64) (records []*model.Record, err error) {
	records = make([]*model.Record, 0, 20)
	err = r.db.Table("last_records").
		Joins("LEFT JOIN `records` r ON r.id = `last_records`.record_id").
		Where("to_id", toId).
		Find(&records).Error
	return records, err
}

func (r *dbRecord) BatchSetRead(ids []int64, toId int64) error {
	return r.db.Model(&model.Record{}).Where("id IN (?) AND to_id = ?", ids, toId).Updates(map[string]interface{}{
		"is_read":    1,
		"updated_at": time.Now().UnixNano(),
	}).Error
}

func (r *dbRecord) PageRecord(page *model.Page, userIdA, userIdB int64, onlyShowNotRead bool) (err error) {
	var queryBuilder = r.db.Model(&model.Record{}).Where("(from_id = ? AND to_id = ?) OR (from_id = ? AND to_id = ?)", userIdA, userIdB, userIdB, userIdA)
	if onlyShowNotRead {
		queryBuilder.Where("is_read", false)
	}
	queryBuilder.Count(&page.Total)
	page.SetPages()
	items := make([]*model.Record, 0, page.Size)
	err = queryBuilder.Order("created_at DESC").Offset(page.Offset()).Limit(page.Size).Find(&items).Error
	page.Records = items
	return
}

func (r *dbRecord) CreateRecord(record *model.Record) (int64, error) {
	record.CreatedAt = time.Now().UnixNano()
	tx := r.db.Begin()
	if err := tx.Create(record).Error; err != nil {
		tx.Rollback()
		return 0, err
	}
	lastRecord := &model.LastRecord{
		FromId:   record.FromId,
		ToId:     record.ToId,
		RecordId: record.Id,
	}
	if err := tx.Create(lastRecord).Error; err != nil {
		tx.Rollback()
		return 0, err
	}
	tx.Commit()
	return record.Id, nil
}

func NewDBRecord(db *gorm.DB) model.RecordStore {
	return &dbRecord{db: db}
}
