package repository

import (
	"github.com/jinzhu/gorm"
	orm "hita/utils/mysql"
)

type Event struct {
	UserId      int64
	TimetableId string `json:"timetableId" gorm:"not null"`
	SubjectId   string `json:"subjectId" gorm:"not null"`
	Id          string `json:"id" gorm:"PRIMARY_KEY"`
	Type        string `json:"type" gorm:"type:enum('CLASS','EXAM','OTHER');default:'OTHER';not null"`
	Name        string `json:"name" gorm:"not null"`
	Place       string `json:"place"`
	Teacher     string `json:"teacher"`
	From        MTime  `json:"from" gorm:"default:null"`
	To          MTime  `json:"to" gorm:"default:null"`
	FromNumber  int8   `json:"fromNumber"`
	LastNumber  int8   `json:"lastNumber"`
}

func (Event) TableName() string {
	return "event"
}

func AddEvents(subjects []Event) {
	for _, tt := range subjects {
		results := orm.DB.Model(Event{}).Where("id = ?", tt.Id).First(&Event{})
		if results.Error != nil {
			if results.Error == gorm.ErrRecordNotFound {
				orm.DB.Model(Event{}).Save(tt)
			}
		} else {
			orm.DB.Model(Event{}).Update(tt)
		}
	}
}

func RemoveEventsInIds(ids []string) {
	orm.DB.Delete(Event{}).Where("id in (?)", ids)
}

func GetEventsInIds(uid string, ids []string) []Event {
	var res []Event
	orm.DB.Raw("select * from event where user_id =? and id in (?)", uid, ids).Scan(&res)
	return res
}
