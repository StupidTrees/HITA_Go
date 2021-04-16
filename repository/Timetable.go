package repository

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	orm "hita/utils/mysql"
	"time"
)

type Timetable struct {
	UserId            int64   `json:"userId" gorm:"column:user_id"`
	Id                string  `json:"id" gorm:"column:id;PRIMARY_KEY"`
	Name              string  `json:"name" gorm:"not null"`
	Code              string  `json:"code" gorm:"not null"`
	StartTime         MTime   `json:"startTime" gorm:"column:start_time;default:null"`
	EndTime           MTime   `json:"endTime" gorm:"column:end_time;default:null"`
	ScheduleStructure SString `json:"scheduleStructure" gorm:"column:schedule_structure"`
}
type MTime time.Time

type SString string

func (t MTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(t).Unix() * 1000)
}

func (s SString) MarshalJSON() ([]byte, error) {
	var l []map[string]map[string]int
	_ = json.Unmarshal([]byte(s), &l)
	return json.Marshal(l)
}

func (Timetable) TableName() string {
	return "timetable"
}

func AddTimetables(timetables []Timetable) {
	for _, tt := range timetables {
		results := orm.DB.Model(Timetable{}).Where("id = ?", tt.Id).First(&Timetable{})
		if results.Error != nil {
			if results.Error == gorm.ErrRecordNotFound {
				orm.DB.Model(Timetable{}).Save(tt)
			}
		} else {
			orm.DB.Model(Timetable{}).Update(tt)
		}
	}
}

func RemoveTimetablesInIds(ids []string) {
	orm.DB.Delete(Timetable{}).Where("id in (?)", ids)
}

func GetTimetablesInIds(uid string, ids []string) []Timetable {
	var res []Timetable
	orm.DB.Raw("select * from timetable where user_id =? and id in (?)", uid, ids).Scan(&res)
	return res
}
