package repository

import (
	orm "hita/utils/mysql"
	"time"
)

type Topic struct {
	Id          int64     `json:"id" gorm:"PRIMARY_KEY;AUTO_INCREMENT"`
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description"`
	Avatar      int64     `json:"avatar"`
	ArticleNum  int       `json:"articleNum" gorm:"not null;default:0"`
	CreateTime  time.Time `gorm:"column:create_time;autoCreateTime:milli"`
}

func (Topic) TableName() string {
	return "topic"
}

func GetHotTopics() (result []Topic, err error) {
	err = orm.DB.Raw("select * from topic order by article_num desc limit 6").Scan(&result).Error
	return
}

func (t *Topic) Get() error {
	return orm.DB.Where("id=?", t.Id).Find(t).Error
}

func SearchTopics(key string, pageSize int, pageNum int) (result []Topic, err error) {
	query := "%" + key + "%"
	err = orm.DB.Raw("select * from topic where name like ? order by article_num desc limit ?,?", query, pageSize*(pageNum-1), pageSize).Scan(&result).Error
	return
}
