package repository

import (
	"hita/utils/logger"
	orm "hita/utils/mysql"
	"time"
)

type Comment struct {
	Id         int64    `json:"id" gorm:"PRIMARY_KEY"`
	Author     User     `gorm:"ForeignKey:AuthorId"`
	AuthorId   int64    `gorm:"not null"`
	Article    Article  `gorm:"ForeignKey:ArticleId"`
	ArticleId  int64    `gorm:"not null"`
	Receiver   User     `gorm:"ForeignKey:ReceiverId"`
	ReceiverId int64    `gorm:"not null"`
	ReplyTo    *Comment `gorm:"ForeignKey:ReplyId"`
	ReplyId    int64
	Context    *Comment `gorm:"ForeignKey:ContextId"`
	ContextId  int64
	Content    string    `gorm:"not null;size:512"`
	LikeNum    int       `gorm:"not null;default:0;size:16"`
	CommentNum int       `gorm:"not null;default:0;size:16"`
	CreateTime time.Time `gorm:"column:create_time;autoCreateTime:milli"`
}

func (Comment) TableName() string {
	return "comment"
}

func (c *Comment) Get() error {
	return orm.DB.Where("id=?", c.Id).Find(c).Error
}
func (c *Comment) Delete() error {
	err := c.Get()
	if err == nil {
		if c.ReplyId > 0 {
			orm.DB.Exec("update comment set comment_num = comment_num-1 where id = ? or id=?", c.ReplyId, c.ContextId)
		}
		orm.DB.Exec("delete from comment where reply_id = ?", c.Id)
		return orm.DB.Where("id=?", c.Id).Delete(c).Error
	} else {
		return err
	}
}

func (a *Comment) Create() error {
	result := orm.DB.Create(a)
	if result.Error != nil {
		logger.Errorln(result.Error)
		return result.Error
	}
	if a.ReplyId > 0 {
		orm.DB.Exec("update comment set comment_num = comment_num+1 where id = ? or id=?", a.ReplyId, a.ContextId)
	}
	return nil
}

func GetCommentsOfArticle(articleId int64, pageSize int, pageNum int) (res []Comment, err error) {
	err = orm.DB.Preload("Author").Preload("Receiver").Preload("ReplyTo").Where("article_id = ? and reply_id = 0", articleId).Limit(pageSize).Offset((pageNum - 1) * pageSize).Order("(like_num+comment_num) DESC").Find(&res).Error
	return
}

func GetCommentsOfComment(commentId int64, pageSize int, pageNum int) (res []Comment, err error) {
	err = orm.DB.Preload("Author").Preload("Receiver").Preload("ReplyTo").Where("context_id = ?", commentId).Limit(pageSize).Offset((pageNum - 1) * pageSize).Order("create_time ASC").Find(&res).Error
	return
}

func GetCommentInfo(commentId int64) (res Comment, err error) {
	err = orm.DB.Preload("Author").Preload("Receiver").Preload("ReplyTo").Where("id = ? ", commentId).First(&res).Error
	return
}
