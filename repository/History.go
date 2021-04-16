package repository

import (
	"encoding/json"
	"errors"
	"github.com/jinzhu/gorm"
	orm "hita/utils/mysql"
)

type History struct {
	Id     int64  `json:"id" gorm:"PRIMARY_KEY"`
	Uid    string `json:"uid" gorm:"column:uid;not null"`
	Table  string `json:"table" gorm:"column:table; not null"`
	Action string `json:"action" gorm:"type:enum('REQUIRE','REMOVE');default:'ADD'"`
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

func GetLatestId(uid string) (id int64, err error) {
	history := History{}
	if orm.DB.Raw("select id from history where uid = ? order by id desc limit 1", uid).Scan(&history).RecordNotFound() {
		err = errors.New("user not exist")
	}
	id = history.Id
	return
}

func SaveHistories(hList []History) {
	for _, h := range hList {
		results := orm.DB.Model(History{}).Where("id = ?", h.Id).First(&History{})
		if results.Error != nil {
			if results.Error == gorm.ErrRecordNotFound {
				orm.DB.Model(History{}).Save(h)
			}
		} else {
			orm.DB.Model(History{}).Update(h)
		}

	}
}

func GetHistoriesAfter(uid string, latestId int64) []History {
	var result []History
	orm.DB.Raw("select * from history where uid = ? and id > ? order by id desc", uid, latestId).Scan(&result)
	return result
}
