package repository

import (
	"gorm.io/gorm/clause"
	orm "hita/utils/mysql"
)

type Event struct {
	UserId      int64
	TimetableId string `json:"timetableId"`
	SubjectId   string `json:"subjectId" `
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
	if len(subjects) > 0 {
		orm.DB.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(&subjects)
	}
}

func ClearEvents(uid int64) {
	orm.DB.Exec("delete from event where user_id=?", uid)
}

func RemoveEventsInIds(ids []string) {
	orm.DB.Where("id in (?)", ids).Delete(Event{})
}

func GetEventsInIds(uid int64, ids []string) []Event {
	var res []Event
	orm.DB.Raw("select * from event where user_id =? and id in (?)", uid, ids).Scan(&res)
	return res
}

func GetEventIds(uid int64) []string {
	var res []string
	orm.DB.Raw("select id from event where user_id =?", uid).Scan(&res)
	return res
}
