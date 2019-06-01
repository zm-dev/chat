package db_store

import (
	"github.com/jinzhu/gorm"
	"github.com/zm-dev/chat_v2/model"
)

type dbUser struct {
	db *gorm.DB
}

func (u *dbUser) UserList(uType model.UserType) (users []*model.User, err error) {
	users = make([]*model.User, 0, 4)
	switch uType {
	case model.TeacherType:
		err = u.db.Model(&model.User{}).Find(&users, map[string]interface{}{"is_teacher": 1}).Error
	case model.StudentType:
		err = u.db.Model(&model.User{}).Find(&users, map[string]interface{}{"is_student": 1}).Error
	case model.AdminType:
		err = u.db.Model(&model.User{}).Find(&users, map[string]interface{}{"is_admin": 1}).Error
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
	return u.db.Omit("created_at").Save(user).Error
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
