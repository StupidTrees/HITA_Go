package repository

import (
	orm "hita/utils/mysql"
	"time"
)

type Message struct {
	Id          int64  `gorm:"PRIMARY_KEY;AUTO_INCREMENT"`
	UserId      int64  `gorm:"not null"`
	OtherId     int64  `gorm:"not null"`
	Action      string `gorm:"type:enum('FOLLOW','UNFOLLOW','COMMENT','LIKE','REPOST');not null; default:'LIKE'"`
	Type        string `gorm:"type:enum('COMMENT','ARTICLE','NONE'); default:'NONE'"`
	Content     string
	ReferenceId string
	Read        bool      `gorm:"not null;default:0"`
	CreateTime  time.Time `gorm:"autoCreateTime:milli"`
}

func (Message) TableName() string {
	return "message"
}

func (m *Message) Create() error {
	return orm.DB.Create(m).Error
}

func CountAllUnread(userId int64) (result int64, err error) {
	err = orm.DB.Raw("select count(*) from message where user_id = ? and `read` = 0 ", userId).Scan(&result).Error
	return
}

func CountUnread(userId int64, tp string) (result int64, err error) {
	err = orm.DB.Raw("select count(*) from message where user_id = ? and `read`= 0 and action = ?", userId, tp).Scan(&result).Error
	return
}

func MarkAllRead(ids []int64) error {
	return orm.DB.Exec("update message set `read` = 1 where id in (?)", ids).Error
}

func GetMessages(userId int64, tp string, pageSize int, pageNum int) (result []Message, err error) {
	err = orm.DB.Raw("select * from message where `action` = ? and user_id=? order by id desc limit ?,?", tp, userId, pageSize*(pageNum-1), pageSize).Scan(&result).Error
	return
}
