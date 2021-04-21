package repository

import (
	orm "hita/utils/mysql"
	"time"
)

type UserLikeArticle struct {
	User       User      `gorm:"ForeignKey:UserId"`
	UserId     int64     `json:"userId" gorm:"not null"`
	Article    Article   `gorm:"ForeignKey:ArticleId"`
	ArticleId  int64     `json:"articleId" gorm:"not null"`
	CreateTime time.Time `gorm:"column:create_time;autoCreateTime:milli"`
}

func (u *UserLikeArticle) Create() error {
	return orm.DB.Create(u).Error
}
func (u *UserLikeArticle) Delete() error {
	return orm.DB.Model(u).Delete(u, "user_id=? and article_id=?", u.UserId, u.ArticleId).Error
}
func (u *UserLikeArticle) GetLikeNum() (a Article, e error) {
	e = orm.DB.Raw("select like_num from article where id = ?", u.ArticleId).Scan(&a).Error
	return
}
func (u *UserLikeArticle) Exist() bool {
	var num int64
	orm.DB.Model(u).Where("user_id=? and article_id=?", u.UserId, u.ArticleId).Limit(1).Count(&num)
	return num != 0
}
