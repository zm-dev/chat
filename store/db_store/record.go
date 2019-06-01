package db_store

import (
	"github.com/jinzhu/gorm"
	"github.com/zm-dev/chat_v2/model"
)

type dbRecord struct {
	db *gorm.DB
}

func (r *dbRecord) BatchSetRead(ids []int64) error {
	return r.db.Model(&model.Record{}).Where("id IN (?)", ids).Update("is_read", 1).Error
}

func (r *dbRecord) PageRecord(page *model.Page, teacherId, studentId int64) (err error) {
	var queryBuilder = r.db.Model(&model.Record{}).Where("teacher_id = ? AND student_id = ?", teacherId, studentId)

	queryBuilder.Count(&page.Total)
	page.SetPages()
	err = queryBuilder.Offset(page.Offset()).Limit(page.Size).Find(&page.Records).Error
	return
}

func (r *dbRecord) CreateRecord(record *model.Record) error {
	return r.db.Create(record).Error;
}

func NewDBRecord(db *gorm.DB) model.RecordStore {
	return &dbRecord{db: db}
}
