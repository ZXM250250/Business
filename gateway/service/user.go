package service

import (
	"errors"
	"gateway/common/response"
	pool "gateway/dao"
	"gateway/model"
	"gateway/tools"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func Login(user *model.User, res *response.UserResponse) (err error) {
	var password = user.Password
	result := pool.GetDB().Where("username=?", user.Username).First(&user)
	if result.Error != nil {
		err = errors.New("未找到用户相关的信息")
		logrus.Info(result.Error)
		return
	}
	if !tools.ComparePasswords(user.Password, password) {
		err = errors.New("登陆失败,密码不正确")
		return
	}
	res.Token, err = tools.GenerateToken(*user)
	res.User = *user
	return
}

func Register(user *model.User, res *response.UserResponse) (err error) {
	salt, err := tools.HashAndSalt(user.Password)
	user.Password = string(salt)
	var userCount int64
	pool.GetDB().Model(user).Where("username=?", user.Username).Count(&userCount)
	if userCount > 0 {
		err = errors.New("用户名已经存在了")
		return
	}

	user.Uuid = uuid.New().String()
	result := pool.GetDB().Create(&user)
	if result.Error != nil {
		err = errors.New("系统内部错误")
		return
	}

	return
}

func ModifyUserAvatar(avatar string, userUuid string) (err error) {
	var queryUser model.User
	pool.GetDB().Where("uuid=?", userUuid).First(&queryUser)
	if tools.NULL_ID == queryUser.Id {
		return errors.New("用户不存在")
	}
	pool.GetDB().Model(&queryUser).Update("avatar", avatar)

	return

}

func ModifyUser(user *model.User) (err error) {
	var queryUser *model.User
	pool.GetDB().First(&queryUser, "username=?,id=?", user.Username, user.Id)
	if tools.NULL_ID == queryUser.Id {
		return errors.New("用户不存在")
	}
	pool.GetDB().Save(user)

	return

}

func GetUserInfo(uuid string) model.User {
	var user model.User
	pool.GetDB().Model(&user).Where("uuid=?", uuid).First(&user)

	return user

}
