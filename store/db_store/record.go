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
	//todo 未完成
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

	lastRecord := &model.LastRecord{}
	userIdA := record.FromId
	userIdB := record.ToId
	r.db.Where("(user_id_a = ? AND user_id_b = ?) OR (user_id_a = ? AND user_id_b = ?)", userIdA, userIdB, userIdB, userIdA).First(&lastRecord)

	tx := r.db.Begin()
	if err := tx.Create(record).Error; err != nil {
		tx.Rollback()
		return 0, err
	}
	// 最新的 recordId
	lastRecord.RecordId = record.Id
	// 如果记录存在则更新
	if lastRecord.UserIdA != 0 && lastRecord.UserIdB != 0 {
		if err := tx.Model(&model.LastRecord{}).
			Where("(user_id_a = ? AND user_id_b = ?) OR (user_id_a = ? AND user_id_b = ?)", userIdA, userIdB, userIdB, userIdA).
			Update("record_id", record.Id).Error; err != nil {
			tx.Rollback()
			return 0, err
		}
	} else {
		lastRecord.UserIdA = userIdA
		lastRecord.UserIdB = userIdB
		if err := tx.Create(lastRecord).Error; err != nil {
			tx.Rollback()
			return 0, err
		}
	}
	tx.Commit()
	return record.Id, nil
}

func NewDBRecord(db *gorm.DB) model.RecordStore {
	return &dbRecord{db: db}
}
