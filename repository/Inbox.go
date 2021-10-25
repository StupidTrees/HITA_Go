package repository

import (
	"hita/utils/logger"
	orm "hita/utils/mysql"
	"time"
)

type Inbox struct {
	Id          int64  `gorm:"PRIMARY_KEY;AUTO_INCREMENT"`
	Content string `json:"content"`
	CreateTime time.Time `gorm:"column:create_time;autoCreateTime:milli"`
}

func (Inbox) TableName() string {
	return "inbox"
}

func (i *Inbox) CreateSuggestion() (err error) {
	result := orm.DB.Create(i)
	if result.Error != nil {
		logger.Errorln(result.Error)
		err = result.Error
		return
	}
	return
}