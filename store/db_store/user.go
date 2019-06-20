package db_store

import (
	"github.com/jinzhu/gorm"
	"github.com/zm-dev/chat/enum"
	"github.com/zm-dev/chat/model"
)

type dbUser struct {
	db *gorm.DB
}

func (u *dbUser) UserList(uType enum.CertificateType, page *model.Page) (err error) {
	switch uType {
	case enum.CertificateTeacher:
		fallthrough
	case enum.CertificateStudent:
		fallthrough
	case enum.CertificateAdmin:
		queryBuilder := u.db.Table("users").Joins("LEFT JOIN `certificates` c ON c.user_id = `users`.id").
			Where("c.type = ?", uType)
		users := make([]*model.User, 0, 10)
		if page.Size+page.Offset() > 0 {
			queryBuilder.Count(&page.Total)
			err = queryBuilder.Offset(page.Offset()).Limit(page.Size).Find(&users).Error
		} else {
			err = queryBuilder.Find(&users).Error
		}
		page.Records = users
		page.SetPages()
	default:
		return model.ErrUserTypeNotExist
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
	return u.db.Model(&model.User{}).Where("id = ?", user.Id).Omit("created_at").Updates(map[string]interface{}{
		"avatar_hash": user.AvatarHash,
		"nick_name":   user.NickName,
		"profile":     user.Profile,
		"gender":      user.Gender,
		"group_id":    user.GroupId,
		"company":     user.Company,
	}).Error
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
