package repository

import (
	"github.com/jinzhu/gorm"
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
	for _, tt := range subjects {
		results := orm.DB.Model(TermSubject{}).Where("id = ?", tt.Id).First(&TermSubject{})
		if results.Error != nil {
			if results.Error == gorm.ErrRecordNotFound {
				orm.DB.Model(TermSubject{}).Save(tt)
			}
		} else {
			orm.DB.Model(TermSubject{}).Update(tt)
		}
	}
}

func RemoveSubjectInIds(ids []string) {
	orm.DB.Delete(TermSubject{}).Where("id in (?)", ids)
}

func GetSubjectsInIds(uid string, ids []string) []TermSubject {
	var res []TermSubject
	orm.DB.Raw("select * from subject where user_id =? and id in (?)", uid, ids).Scan(&res)
	return res
}
