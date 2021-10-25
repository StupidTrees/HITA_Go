package repository

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"hita/utils"
	"hita/utils/logger"
	orm "hita/utils/mysql"
	"time"
)

type Article struct {
	Id         int64 `json:"id" gorm:"PRIMARY_KEY"`
	Author     User  `gorm:"ForeignKey:AuthorId"`
	AuthorId   int64 `gorm:"not null"`
	Topic      Topic `json:"topic" gorm:"ForeignKey:TopicId"`
	TopicId    sql.NullInt64
	RepostFrom *Article `gorm:"ForeignKey:RepostId"`
	RepostId   int64
	Content    string    `gorm:"not null;size:512"`
	LikeNum    int       `gorm:"not null;default:0;size:16"`
	Images     MIntArray `gorm:"column:images"`
	CommentNum int       `gorm:"not null;default:0;size:16"`
	Type       string    `json:"type" gorm:"type:enum('NORMAL','VOTE');default:'NORMAL';not null"`
	UpNum      int       `gorm:"default:0;not null"`
	DownNum    int       `gorm:"default:0;not null"`
	Anonymous  bool
	CreateTime time.Time `gorm:"column:create_time;autoCreateTime:milli"`
	UpdateTime int64     `gorm:"column:update_time;autoUpdateTime:milli"`
}

type MInt64 int64
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
	err := orm.DB.Preload("Author").Preload("Topic").Where("id=?", a.Id).Find(a).Error
	a.eraseName()
	return err
}
func (a *Article) Delete() error {
	return orm.DB.Where("id=?", a.Id).Delete(a).Error
}

func GetFollowingPosts(userId int64, beforeTime utils.Long, afterTime utils.Long, pageSize int) (res []Article, err error) {
	beforeTs := beforeTime.ToTime().UTC()
	afterTs := afterTime.ToTime().UTC()
	err = orm.DB.Preload("Author").Preload("Topic").Preload("RepostFrom").Where("create_time < ? and  create_time>? and ( author_id = ? or author_id in (?) )", beforeTs, afterTs, userId,
		orm.DB.Raw("select u.id from user as u where exists(select * from follows where user_id = ? and following_id = u.id)", userId)).Order("id DESC").Limit(pageSize).Find(&res).Error
	filter(res)
	return
}

func  GetUsersPosts(userId int64, beforeTime utils.Long, afterTime utils.Long, pageSize int) (res []Article, err error) {
	beforeTs := beforeTime.ToTime().UTC()
	afterTs := afterTime.ToTime().UTC()
	err = orm.DB.Preload("Author").Preload("Topic").Preload("RepostFrom").Where("create_time < ? and  create_time>? and ( author_id = ? )", beforeTs, afterTs, userId).Order("id DESC").Limit(pageSize).Find(&res).Error
	filter(res)
	return
}

func GetStarredPosts(userId int64, beforeTime utils.Long, afterTime utils.Long, pageSize int) (res []Article, err error) {
	beforeTs := beforeTime.ToTime().UTC()
	afterTs := afterTime.ToTime().UTC()
	err = orm.DB.Preload("Author").Preload("Topic").Preload("RepostFrom").Where("create_time < ? and  create_time>? and id in (?)", beforeTs, afterTs,
		orm.DB.Raw("select article_id from stars where user_id = ?", userId)).Order("id DESC").Limit(pageSize).Find(&res).Error
	filter(res)
	return
}

func GetTopicPosts(topicId int64, beforeTime utils.Long, afterTime utils.Long, pageSize int) (res []Article, err error) {
	beforeTs := beforeTime.ToTime().UTC()
	afterTs := afterTime.ToTime().UTC()
	err = orm.DB.Preload("Author").Preload("Topic").Preload("RepostFrom").Where("create_time < ? and  create_time>? and topic_id =?", beforeTs, afterTs, topicId).Order("id DESC").Limit(pageSize).Find(&res).Error
	filter(res)
	return
}

func GetReposts(articleId int64, beforeTime utils.Long, afterTime utils.Long, pageSize int) (res []Article, err error) {
	beforeTs := beforeTime.ToTime().UTC()
	afterTs := afterTime.ToTime().UTC()
	err = orm.DB.Preload("Author").Preload("Topic").Preload("RepostFrom").Where("create_time < ? and  create_time>? and ( repost_id = ? )", beforeTs, afterTs, articleId).Order("id DESC").Limit(pageSize).Find(&res).Error
	filter(res)
	return
}

func GetAllPosts(beforeTime utils.Long, afterTime utils.Long, pageSize int) (res []Article, err error) {
	err = orm.DB.Preload("Author").Preload("Topic").Preload("RepostFrom").Where("create_time < ? and  create_time>? ", beforeTime.ToTime().UTC(), afterTime.ToTime().UTC()).Order("id DESC").Limit(pageSize).Find(&res).Error
	filter(res)
	return
}

func SearchPosts(beforeTime utils.Long, afterTime utils.Long, pageSize int, extra string) (res []Article, err error) {
	if len([]rune(extra)) == 0 {
		return []Article{}, nil
	}
	likeS := "%" + extra + "%"
	err = orm.DB.Preload("Author").Preload("Topic").Preload("RepostFrom").Where("content like ? and create_time < ? and  create_time>? ", likeS, beforeTime.ToTime().UTC(), afterTime.ToTime().UTC()).Order("id DESC").Limit(pageSize).Find(&res).Error
	filter(res)
	return
}

func filter(articles []Article) {
	for index, a := range articles {
		a.eraseName()
		articles[index] = a
	}
}

var  AnonymousId int64= 0
func (a *Article) eraseName() {
	if a.Anonymous {
		a.Author = User{
			Id:       AnonymousId,
			Nickname: "匿名用户",
			Avatar:   3,
		}
	}
}
