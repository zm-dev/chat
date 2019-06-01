package model

type Record struct {
	Id        int64  `gorm:"type:BIGINT AUTO_INCREMENT;PRIMARY_KEY;NOT NUll" json:"id"`
	TeacherId int64  `gorm:"type:BIGINT;NOT NUll" json:"teacher_id"`    // 教师ID
	StudentId int64  `gorm:"type:BIGINT;NOT NUll" json:"student_id"`    // 学生ID
	From      int8   `gorm:"type:TINYINT;NOT NUll" json:"from"`         // 记录来源（教师、学生）
	Content   string `gorm:"type:varchar(512);NOT NULL" json:"content"` // 正文
	IsRead    bool   `gorm:"type:TINYINT" json:"is_read"`               // 是否已经阅读
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

type RecordStore interface {
	// 批量设置聊天记录为已读状态
	BatchSetRead(ids []int64) error
	// 分页获取聊天记录列表，按照创建时间倒序排序
	PageRecord(page *Page, teacherId, studentId int64) (err error)
	// 创建一条聊天记录
	CreateRecord(record *Record) error
}

type RecordService interface {
	RecordStore
}
