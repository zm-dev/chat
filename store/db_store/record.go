package db_store

import (
	"github.com/jinzhu/gorm"
	"github.com/zm-dev/chat/model"
	"time"
)

type dbRecord struct {
	db *gorm.DB
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
	err = queryBuilder.Offset(page.Offset()).Limit(page.Size).Find(&items).Error
	page.Records = items
	return
}

func (r *dbRecord) CreateRecord(record *model.Record) error {
	record.CreatedAt = time.Now().UnixNano()
	return r.db.Create(record).Error
}

func NewDBRecord(db *gorm.DB) model.RecordStore {
	return &dbRecord{db: db}
}
