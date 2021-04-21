package repository

import (
	"gorm.io/gorm/clause"
	orm "hita/utils/mysql"
)

type TermSubject struct {
	UserId      int64   `json:"userId" gorm:"not null"`
	TimetableId string  `json:"timetableId" gorm:"not null"`
	Id          string  `json:"id" gorm:"PRIMARY_KEY"`
	Type        string  `json:"type" gorm:"type:enum('COM_A', 'COM_B', 'OPT_A', 'OPT_B', 'MOOC');default:'COM_B';not null"`
	Name        string  `json:"name" gorm:"not null"`
	Field       string  `json:"field"`
	School      string  `json:"school"`
	Code        string  `json:"code"`
	Key         string  `json:"key"`
	Color       int32   `json:"color"`
	Credit      float32 `json:"credit"`
	CountInSPA  bool    `json:"countInSpa"`
}

func (TermSubject) TableName() string {
	return "subject"
}

func AddSubjects(subjects []TermSubject) {
	if len(subjects) > 0 {
		orm.DB.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(&subjects)
	}

}

func RemoveSubjectInIds(ids []string) {
	orm.DB.Where("id in (?)", ids).Delete(TermSubject{})
}

func GetSubjectsInIds(uid string, ids []string) []TermSubject {
	var res []TermSubject
	orm.DB.Raw("select * from subject where user_id =? and id in (?)", uid, ids).Scan(&res)
	return res
}
