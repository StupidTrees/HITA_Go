package repository

import (
	"errors"
	orm "hita/utils/mysql"
	"time"
)

type Vote struct {
	User       User      `gorm:"ForeignKey:UserId"`
	UserId     int64     `json:"userId" gorm:"not null"`
	Article    Article   `gorm:"ForeignKey:ArticleId"`
	ArticleId  int64     `json:"articleId" gorm:"not null"`
	Up         bool      `gorm:"not null;default:false"`
	CreateTime time.Time `gorm:"column:create_time;autoCreateTime:milli"`
}

func (Vote) TableName() string {
	return "vote"
}

func (u *Vote) Create() error {
	return orm.DB.Create(u).Error
}

func (u *Vote) GetVoteNum() (a Article, e error) {
	e = orm.DB.Raw("select up_num,down_num from article where id = ?", u.ArticleId).Scan(&a).Error
	return
}

func (u *Vote) Find() (err error) {
	var votes []Vote
	err = orm.DB.Model(u).Raw("select * from vote where user_id=? and article_id=?", u.UserId, u.ArticleId).Scan(&votes).Error
	if len(votes)>0{
		*u = votes[0]
		return nil
	}else{
		return errors.New("not voted")
	}
}


