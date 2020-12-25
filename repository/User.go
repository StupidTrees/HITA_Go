package repository

import (
	"hita/lib/logger"
	orm "hita/lib/mysql"
	"time"
)

type User struct {
	Id         int64     `json:"id" gorm:"PRIMARY_KEY"`
	UserName   string    `json:"username" gorm:"column:username; unique_index:username_idx; not null"`
	Password   string    `json:"password" gorm:"column:password; not null"`
	Nickname   string    `json:"nickname" gorm:"column:nickname"`
	Gender     string    `json:"gender" gorm:"type:enum('OTHER','MALE','FEMALE');default:OTHER"`
	StudentId  string    `json:"student_id"`
	School     string    `json:"school"`
	CreateTime time.Time `json:"create_time" gorm:"column:createtime;default:null"`
	UpdateTime time.Time `json:"update_time" gorm:"column:updatetime;default:null"`
}

func (User) TableName() string {
	return "user"
}

func (user *User) AddUser() (id int64, err error) {

	result := orm.DB.Create(user)
	id = user.Id
	if result.Error != nil {
		logger.Errorln(result.Error)
		err = result.Error
		return
	}
	return

}

//func (p *Person) GetPersons() (persons []Person, err error) {
//	persons = make([]Person, 0)
//	rows, err := db.SqlDB.Query("SELECT id, first_name, last_name FROM person")
//	defer rows.Close()
//	if err != nil {
//		return
//	}
//	for rows.Next() {
//		var person Person
//		rows.Scan(&person.Id, &person.FirstName, &person.LastName)
//		persons = append(persons, person)
//	}
//	if err = rows.Err(); err != nil {
//		return
//	}
//	return
//}
//
//func (p *Person) GetPerson() (person Person, err error) {
//	err = db.SqlDB.QueryRow("SELECT id, first_name, last_name FROM person WHERE id=?", p.Id).Scan(
//		&person.Id, &person.FirstName, &person.LastName,
//	)
//	return
//}
//
//func (p *Person) ModPerson() (ra int64, err error) {
//	stmt, err := db.SqlDB.Prepare("UPDATE person SET first_name=?, last_name=? WHERE id=?")
//	defer stmt.Close()
//	if err != nil {
//		return
//	}
//	rs, err := stmt.Exec(p.FirstName, p.LastName, p.Id)
//	if err != nil {
//		return
//	}
//	ra, err = rs.RowsAffected()
//	return
//}
//
//func (p *Person) DelPerson() (ra int64, err error) {
//	rs, err := db.SqlDB.Exec("DELETE FROM person WHERE id=?", p.Id)
//	if err != nil {
//		log.Fatalln(err)
//	}
//	ra, err = rs.RowsAffected()
//	return
//}
