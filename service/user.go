package service

import (
	"context"
	"github.com/zm-dev/chat/enum"
	"github.com/zm-dev/chat/errors"
	"github.com/zm-dev/chat/model"
	"github.com/zm-dev/chat/pkg/hasher"
)

type userService struct {
	model.UserStore
	cs   model.CertificateStore
	tSvc model.TicketService
	h    hasher.Hasher
}

func (uSvc *userService) TeacherList() ([]*model.User, error) {
	return uSvc.UserStore.UserList(enum.TeacherType)
}

func (uSvc *userService) StudentList() ([]*model.User, error) {
	return uSvc.UserStore.UserList(enum.StudentType)
}

func (uSvc *userService) UserLogin(account, password string) (ticket *model.Ticket, err error) {
	c, err := uSvc.cs.CertificateLoadByAccount(account)
	if err != nil {
		if uSvc.cs.CertificateIsNotExistErr(err) { //账号不存在
			err = errors.ErrAccountNotFound()
		}
		return nil, err
	}
	user, err := uSvc.UserStore.UserLoad(c.UserId)
	if err != nil {
		return nil, err
	}
	if uSvc.h.Check(password, user.Password) {
		// 登录成功
		return uSvc.tSvc.TicketGen(user.Id)
	}

	return nil, errors.ErrPassword()
}

func (uSvc *userService) UserRegister(account string, certificateType model.CertificateType, password string) (userId int64, err error) {
	if exist, err := uSvc.cs.CertificateExist(account); err != nil {
		return 0, err
	} else if exist {
		return 0, errors.ErrAccountAlreadyExisted()
	}
	user := &model.User{
		Password: uSvc.h.Make(password),
		PwPlain:  password,
		Gender:   enum.GenderSecrecy,
		// 学生注册
		IsStudent: true,
	}
	if err := uSvc.UserStore.UserCreate(user); err != nil {
		return 0, err
	}
	certificate := &model.Certificate{UserId: user.Id, Account: account, Type: certificateType}
	if err := uSvc.cs.CertificateCreate(certificate); err != nil {
		return 0, err
	}
	return user.Id, nil
}

func (uSvc *userService) UserUpdatePassword(userId int64, newPassword string) error {
	return uSvc.UserStore.UserUpdate(&model.User{
		Password: uSvc.h.Make(newPassword),
		PwPlain:  newPassword,
	})
}

func UserLoad(ctx context.Context, id int64) (*model.User, error) {
	return FromContext(ctx).UserLoad(id)
}

func UserLogin(ctx context.Context, account, password string) (*model.Ticket, error) {
	return FromContext(ctx).UserLogin(account, password)
}

func UserRegister(ctx context.Context, account string, certificateType model.CertificateType, password string) (userId int64, err error) {
	return FromContext(ctx).UserRegister(account, certificateType, password)
}

func UserUpdatePassword(ctx context.Context, userId int64, newPassword string) error {
	return FromContext(ctx).UserUpdatePassword(userId, newPassword)
}

func TeacherList(ctx context.Context) ([]*model.User, error) {
	return FromContext(ctx).TeacherList()
}

func w(ctx context.Context) ([]*model.User, error) {
	return FromContext(ctx).StudentList()
}

func NewUserService(us model.UserStore, cs model.CertificateStore, tSvc model.TicketService, h hasher.Hasher) model.UserService {
	return &userService{us, cs, tSvc, h}
}
