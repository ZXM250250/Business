package dao

import (
	"gorm.io/gorm"
)

type GrpcRule struct {
	ID             int64  `json:"id" gorm:"primary_key"`
	ServiceID      int64  `json:"service_id" gorm:"column:service_id" description:"服务id	"`
	Port           int    `json:"port" gorm:"column:port" description:"端口	"`
	HeaderTransfor string `json:"header_transfor" gorm:"column:header_transfor" description:"header转换支持增加(add)、删除(del)、修改(edit) 格式: add headname headvalue"`
}

func (t *GrpcRule) TableName() string {
	return "gateway_service_grpc_rule"
}

func (t *GrpcRule) Find(search *GrpcRule) (*GrpcRule, error) {
	model := &GrpcRule{}
	err := GetDB().Where(search).Find(model).Error
	return model, err
}

func (t *GrpcRule) Save() error {
	if err := GetDB().Save(t).Error; err != nil {
		return err
	}
	return nil
}

func (t *GrpcRule) ListByServiceID(serviceID int64) ([]GrpcRule, int64, error) {
	var list []GrpcRule
	var count int64
	query := GetDB()
	query = query.Table(t.TableName()).Select("*")
	query = query.Where("service_id=?", serviceID)
	err := query.Order("id desc").Find(&list).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, 0, err
	}
	errCount := query.Count(&count).Error
	if errCount != nil {
		return nil, 0, err
	}
	return list, count, nil
}
