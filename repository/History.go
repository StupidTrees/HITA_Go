package repository

import (
	"encoding/json"
	"gorm.io/gorm/clause"
	orm "hita/utils/mysql"
)

type History struct {
	Id     int64  `json:"id" gorm:"PRIMARY_KEY"`
	Uid    string `json:"uid" gorm:"column:uid;not null"`
	Table  string `json:"table" gorm:"column:table; not null"`
	Action string `json:"action" gorm:"type:enum('REQUIRE','REMOVE','CLEAR');default:'REQUIRE'"`
	Ids    string `json:"ids"`
}

func (History) TableName() string {
	return "history"
}

func (h History) GetIds() []string {
	var list []string
	_ = json.Unmarshal([]byte(h.Ids), &list)
	return list
}

func (h *History) SetIds(list []string) {
	str, _ := json.Marshal(list)
	h.Ids = string(str)
}

func GetLatestId(uid int64) (id int64, err error) {
	history := History{}
	orm.DB.Raw("select id from history where uid = ? order by id desc limit 1", uid).Scan(&history)
	id = history.Id
	return
}

func SaveHistories(hList []History) {
	if len(hList) > 0 {
		orm.DB.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(&hList)
	}
}

func ClearHistories(uid int64, reserve int64) {
	orm.DB.Exec("delete from history where uid = ? and id <> ?", uid, reserve)
}
func GetHistoriesAfter(uid string, latestId int64) []History {
	var result []History
	orm.DB.Raw("select * from history where uid = ? and id > ? order by id desc", uid, latestId).Scan(&result)
	return result
}
