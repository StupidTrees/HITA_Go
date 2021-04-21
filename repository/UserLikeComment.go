package repository

import (
	orm "hita/utils/mysql"
	"time"
)

type UserLikeComment struct {
	User       User      `gorm:"ForeignKey:UserId"`
	UserId     int64     `json:"userId" gorm:"not null"`
	Comment    Comment   `gorm:"ForeignKey:CommentId"`
	CommentId  int64     `json:"commentId" gorm:"not null"`
	CreateTime time.Time `gorm:"column:create_time;autoCreateTime:milli"`
}

func (u *UserLikeComment) Create() error {
	return orm.DB.Create(u).Error
}
func (u *UserLikeComment) Delete() error {
	return orm.DB.Model(u).Delete(u, "user_id=? and comment_id=?", u.UserId, u.CommentId).Error
}
func (u *UserLikeComment) GetLikeNum() (a Comment, e error) {
	e = orm.DB.Raw("select like_num from comment where id = ?", u.CommentId).Scan(&a).Error
	return
}
func (u *UserLikeComment) Exist() bool {
	var num int64
	orm.DB.Model(u).Where("user_id=? and comment_id=?", u.UserId, u.CommentId).Limit(1).Count(&num)
	return num != 0
}
