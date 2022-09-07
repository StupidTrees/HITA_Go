package repository

import (
	"database/sql/driver"
	"encoding/json"
	"gorm.io/gorm/clause"
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
func (t *MTime) UnmarshalJSON(b []byte) error {
	var i int64
	err := json.Unmarshal(b, &i)
	if err == nil {
		*t = MTime(time.Unix(i/1000, 0))
	}
	return err
}

// 写入数据库之前，对数据做类型转换
func (t MTime) Value() (driver.Value, error) {
	x := time.Time(t)
	return x, nil
}

// 将数据库中取出的数据，赋值给目标类型
func (t *MTime) Scan(v interface{}) error {
	x := v.(time.Time)
	*t = MTime(x)
	return nil
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
	if len(timetables) > 0 {
		orm.DB.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(&timetables)
	}

}

func ClearTimetables(uid int64) {
	orm.DB.Exec("delete from timetable where user_id=?", uid)
}

func RemoveTimetablesInIds(ids []string) {
	orm.DB.Where("id in (?)", ids).Delete(Timetable{})
}

func GetTimetablesInIds(uid int64, ids []string) []Timetable {
	var res []Timetable
	orm.DB.Raw("select * from timetable where user_id =? and id in (?)", uid, ids).Scan(&res)
	return res
}

func GetTimetableIds(uid int64) []string {
	var res []string
	orm.DB.Raw("select id from timetable where user_id =?", uid).Scan(&res)
	return res
}
