package repository

import (
	orm "hita/utils/mysql"
	"time"
)

type Follow struct {
	User        User      `gorm:"ForeignKey:UserId"`
	UserId      int64     `json:"userId" gorm:"not null"`
	Following   User      `gorm:"ForeignKey:FollowingId"`
	FollowingId int64     `json:"followingId" gorm:"not null"`
	CreateTime  time.Time `gorm:"column:create_time;autoCreateTime:milli"`
}

func (u *Follow) Create() error {
	return orm.DB.Create(u).Error
}
func (u *Follow) Delete() error {
	return orm.DB.Model(u).Delete(u, "user_id=? and following_id=?", u.UserId, u.FollowingId).Error
}

func (u *Follow) Exist() bool {
	var num int64
	orm.DB.Model(u).Where("user_id = ? and  following_id = ?", u.UserId, u.FollowingId).Limit(1).Count(&num)
	return num != 0
}
