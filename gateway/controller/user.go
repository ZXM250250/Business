package controller

import (
	"gateway/common/response"
	"gateway/config"
	"gateway/model"
	"gateway/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"net/http"
	"strings"
)

func Login(c *gin.Context) {
	var user model.User
	var res response.UserResponse
	c.ShouldBind(&user)
	if err := service.Login(&user, &res); err != nil {
		c.JSON(http.StatusOK, response.FailMsg(err.Error()))
		logrus.WithField("login", "user").Info(err)
		return
	}
	c.JSON(http.StatusOK, response.SuccessMsg(res))

}

func Register(c *gin.Context) {
	var user model.User
	var res response.UserResponse
	c.ShouldBind(&user)
	if err := service.Register(&user, &res); err != nil {
		c.JSON(http.StatusOK, response.FailMsg(err.Error()))
		logrus.WithField("login", "user").Info(err)

		return
	}
	c.JSON(http.StatusOK, response.SuccessMsg(res))
}

func ModifyUserAvatar(c *gin.Context) {
	namePrefix := uuid.New().String()
	userUuid := c.PostForm("uuid")
	file, _ := c.FormFile("file")
	fileName := file.Filename
	index := strings.LastIndex(fileName, ".")
	suffix := fileName[index:]
	newFileName := config.GetConfig().StaticPath.FilePath + namePrefix + suffix
	err := c.SaveUploadedFile(file, newFileName)
	if err != nil {
		logrus.WithField("login", "user").Info(err)
		return
	}

	logrus.Info("fileName=?", newFileName)
	logrus.Info("userUuid=?", userUuid)
	err = service.ModifyUserAvatar(newFileName, userUuid)
	if err != nil {
		c.JSON(http.StatusOK, response.FailMsg(err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.SuccessMsg(newFileName))

}

func ModifyUser(c *gin.Context) {
	var user model.User
	c.ShouldBind(&user)
	if err := service.ModifyUser(&user); err != nil {
		c.JSON(http.StatusOK, response.FailMsg(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessMsg(user))

}

func GetUserInfo(c *gin.Context) {
	uuid := c.Param("uuid")
	c.JSON(http.StatusOK, response.SuccessMsg(service.GetUserInfo(uuid)))
}
