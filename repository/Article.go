package repository

import (
	"database/sql/driver"
	"encoding/json"
	"hita/utils"
	"hita/utils/logger"
	orm "hita/utils/mysql"
	"time"
)

type Article struct {
	Id         int64    `json:"id" gorm:"PRIMARY_KEY"`
	Author     User     `gorm:"ForeignKey:AuthorId"`
	AuthorId   int64    `gorm:"not null"`
	RepostFrom *Article `gorm:"ForeignKey:RepostId"`
	RepostId   int64
	Content    string    `gorm:"not null;size:512"`
	LikeNum    int       `gorm:"not null;default:0;size:16"`
	Images     MIntArray `gorm:"column:images"`
	CommentNum int       `gorm:"not null;default:0;size:16"`
	CreateTime time.Time `gorm:"column:create_time;autoCreateTime:milli"`
	UpdateTime int64     `gorm:"column:update_time;autoUpdateTime:milli"`
}

type MIntArray []int64

func (Article) TableName() string {
	return "article"
}

func (t MIntArray) MarshalJSON() ([]byte, error) {
	return json.Marshal([]int64(t))
}
func (t *MIntArray) UnmarshalJSON(b []byte) error {
	var i []int64
	err := json.Unmarshal(b, &i)
	if err == nil {
		*t = i
	}
	return err
}

// 写入数据库之前，对数据做类型转换
func (t MIntArray) Value() (driver.Value, error) {
	js, err := t.MarshalJSON()
	return string(js), err
}

// 将数据库中取出的数据，赋值给目标类型
func (t *MIntArray) Scan(v interface{}) error {
	x := v.([]byte)
	return t.UnmarshalJSON(x)
}

func (a *Article) Create() error {
	result := orm.DB.Create(a)
	if result.Error != nil {
		logger.Errorln(result.Error)
		return result.Error
	}
	return nil
}
func (a *Article) Get() error {
	return orm.DB.Preload("Author").Where("id=?", a.Id).Find(a).Error
}
func (a *Article) Delete() error {
	return orm.DB.Where("id=?", a.Id).Delete(a).Error
}

func GetFollowingPosts(userId int64, beforeTime utils.Long, afterTime utils.Long, pageSize int) (res []Article, err error) {
	beforeTs := beforeTime.ToTime().UTC()
	afterTs := afterTime.ToTime().UTC()
	err = orm.DB.Preload("Author").Preload("RepostFrom").Where("create_time < ? and  create_time>? and ( author_id = ? or author_id in (?) )", beforeTs, afterTs, userId,
		orm.DB.Raw("select u.id from user as u where exists(select * from follows where user_id = ? and following_id = u.id)", userId)).Order("id DESC").Limit(pageSize).Find(&res).Error
	return
}

func GetUsersPosts(userId int64, beforeTime utils.Long, afterTime utils.Long, pageSize int) (res []Article, err error) {
	beforeTs := beforeTime.ToTime().UTC()
	afterTs := afterTime.ToTime().UTC()
	err = orm.DB.Preload("Author").Preload("RepostFrom").Where("create_time < ? and  create_time>? and ( author_id = ? )", beforeTs, afterTs, userId).Order("id DESC").Limit(pageSize).Find(&res).Error
	return
}

func GetReposts(articleId int64, beforeTime utils.Long, afterTime utils.Long, pageSize int) (res []Article, err error) {
	beforeTs := beforeTime.ToTime().UTC()
	afterTs := afterTime.ToTime().UTC()
	err = orm.DB.Preload("Author").Preload("RepostFrom").Where("create_time < ? and  create_time>? and ( repost_id = ? )", beforeTs, afterTs, articleId).Order("id DESC").Limit(pageSize).Find(&res).Error
	return
}

func GetAllPosts(beforeTime utils.Long, afterTime utils.Long, pageSize int) (res []Article, err error) {
	err = orm.DB.Preload("Author").Preload("RepostFrom").Where("create_time < ? and  create_time>? ", beforeTime.ToTime().UTC(), afterTime.ToTime().UTC()).Order("id DESC").Limit(pageSize).Find(&res).Error
	return
}

func SearchPosts(beforeTime utils.Long, afterTime utils.Long, pageSize int, extra string) (res []Article, err error) {
	if len([]rune(extra)) == 0 {
		return []Article{}, nil
	}
	likeS := "%" + extra + "%"
	err = orm.DB.Preload("Author").Preload("RepostFrom").Where("content like ? and create_time < ? and  create_time>? ", likeS, beforeTime.ToTime().UTC(), afterTime.ToTime().UTC()).Order("id DESC").Limit(pageSize).Find(&res).Error
	return
}
