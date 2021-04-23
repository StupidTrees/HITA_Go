package repository

import (
	"hita/utils/logger"
	orm "hita/utils/mysql"
)

type Image struct {
	Id       int64  `json:"id" gorm:"PRIMARY_KEY;AUTO_INCREMENT"`
	Filename string `json:"filename"`
	Type     string `json:"type" gorm:"type:enum('AVATAR','POST','OTHER');default:'OTHER';not null"`
}

func (Image) TableName() string {
	return "image"
}

func (i *Image) Create() error {
	result := orm.DB.Create(i)
	if result.Error != nil {
		logger.Errorln(result.Error)
		return result.Error
	}
	return nil
}

func (i *Image) Delete() error {
	return orm.DB.Where("id=?", i.Id).Delete(i).Error
}
func (i *Image) Find() error {
	return orm.DB.Where("id=?", i.Id).Find(i).Error
}
