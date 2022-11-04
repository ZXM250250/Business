package dao

import (
	"fmt"
	"gateway/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var _db *gorm.DB

func init() {
	username := config.GetConfig().Sql.User
	password := config.GetConfig().Sql.Password
	host := config.GetConfig().Sql.Host //数据库地址，可以是Ip或者域名
	port := config.GetConfig().Sql.Port //数据库端口
	Dbname := config.GetConfig().Sql.Name
	timeout := config.GetConfig().Sql.TimeOut
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local&timeout=%s", username, password, host, port, Dbname, timeout)
	var err error
	fmt.Println(dsn)
	_db, err = gorm.Open(mysql.New(mysql.Config{
		DSN:               dsn,
		DefaultStringSize: 256, //string 类型的默认长度

	}), &gorm.Config{})
	if err != nil {
		panic("连接数据库失败, error=" + err.Error())
	}
	sqlDB, _ := _db.DB()
	err = sqlDB.Ping()
	if err != nil {
		panic("数据库不通" + err.Error())
	}
	//设置数据库连接池参数
	sqlDB.SetMaxOpenConns(100) //设置数据库连接池最大连接数
	sqlDB.SetMaxIdleConns(20)  //连接池最大允许的空闲连接数，如果没有sql任务需要执行的连接数大于20，超过的连接会被连接池关闭。
}

func GetDB() *gorm.DB {
	return _db
}
