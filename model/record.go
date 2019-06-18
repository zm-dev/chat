package model

import (
	"errors"
	"time"
)

type Record struct {
	Id        int64     `gorm:"type:BIGINT AUTO_INCREMENT;PRIMARY_KEY;NOT NUll" json:"id"`
	FromId    int64     `gorm:"type:BIGINT;NOT NUll" json:"from_id"`
	ToId      int64     `gorm:"type:BIGINT;NOT NUll" json:"to_id"`
	Msg       string    `gorm:"type:varchar(512);NOT NULL" json:"msg"` // 正文
	IsRead    bool      `gorm:"type:TINYINT" json:"is_read"`           // 是否已经阅读
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type LastRecord struct {
	UserIdA  int64 `gorm:"type:BIGINT;NOT NUll;unique_index:user_id_a_user_id_b_ux" json:"user_id_a"`
	UserIdB  int64 `gorm:"type:BIGINT;NOT NUll;unique_index:user_id_a_user_id_b_ux" json:"user_id_b"`
	RecordId int64 `gorm:"type:BIGINT;NOT NUll" json:"record_id"`
}
 var  ErrFromUserEqualToUser = errors.New("不允许自己和自己聊天")

type RecordStore interface {
	// 批量设置聊天记录为已读状态
	BatchSetRead(ids []int64, toId int64) error
	// 分页获取聊天记录列表，按照创建时间倒序排序
	PageRecord(page *Page, userIdA, userIdB int64, onlyShowNotRead bool) (err error)
	// 创建一条聊天记录
	CreateRecord(record *Record) (int64, error)
	// 最近的聊天记录(userId 必须传自己的ID，获取和自己有关的消息)
	LastRecordList(userIdA int64, size int) (records []*Record, err error)
	// 未读的消息数量
	GetNotReadRecordCount(fromId, toId int64) int32
}

type RecordService interface {
	RecordStore
}
