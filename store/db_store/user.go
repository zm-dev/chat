package db_store

import (
	"github.com/jinzhu/gorm"
	"github.com/zm-dev/chat/enum"
	"github.com/zm-dev/chat/model"
)

type dbUser struct {
	db *gorm.DB
}

func (u *dbUser) UserList(uType enum.CertificateType) (users []*model.User, err error) {
	users = make([]*model.User, 0, 4)
	switch uType {
	case enum.CertificateTeacher:
		fallthrough
	case enum.CertificateStudent:
		fallthrough
	case enum.CertificateAdmin:
		err = u.db.Table("users").Joins("LEFT JOIN `certificates` c ON c.user_id = `users`.id").
			Where("c.type = ?", uType).Find(&users).Error
	default:
		return nil, model.ErrUserTypeNotExist
	}
	return
}

func (u *dbUser) UserExist(id int64) (bool, error) {
	var count uint8
	err := u.db.Model(model.User{}).Where(model.User{Id: id}).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (u *dbUser) UserLoad(id int64) (user *model.User, err error) {
	if id <= 0 {
		return nil, model.ErrUserNotExist
	}
	user = &model.User{}
	err = u.db.Where(model.User{Id: id}).First(user).Error
	if gorm.IsRecordNotFoundError(err) {
		err = model.ErrUserNotExist
	}
	return
}

func (u *dbUser) UserUpdate(user *model.User) error {
	if user.Id <= 0 {
		return model.ErrUserNotExist
	}
	return u.db.Model(&model.User{}).Omit("created_at").Updates(user).Error
}

func (u *dbUser) UserCreate(user *model.User) error {
	return u.db.Create(user).Error
}

func (u *dbUser) UserListByUserIds(userIds []interface{}) (users []*model.User, err error) {
	if len(userIds) == 0 {
		return
	}
	users = make([]*model.User, 0, len(userIds))
	err = u.db.Where("id in (?)", userIds).Find(&users).Error
	return
}

func NewDBUser(db *gorm.DB) model.UserStore {
	return &dbUser{db: db}
}
