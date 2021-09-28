package repository

import (
	orm "hita/utils/mysql"
)

type Info struct {
	Key   string `json:"key" gorm:"PRIMARY_KEY"`
	Value string `json:"value"`
}
func (Info) TableName() string {
	return "info"
}
func (u *Info) Get() (err error) {
	err = orm.DB.Where("`key` = ?", u.Key).First(u).Error
	return err
}
