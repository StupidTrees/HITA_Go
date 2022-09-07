package repository

import (
	"errors"
	"gorm.io/gorm"
	"hita/utils/logger"
	orm "hita/utils/mysql"
	"time"
)

type User struct {
	Id           int64     `json:"id" gorm:"PRIMARY_KEY;AUTO_INCREMENT"`
	UserName     string    `json:"username" gorm:"column:username; unique_index:username_idx; not null"`
	Password     string    `json:"password" gorm:"column:password; not null"`
	Nickname     string    `json:"nickname" gorm:"column:nickname"`
	Gender       string    `json:"gender" gorm:"type:enum('OTHER','MALE','FEMALE');default:OTHER"`
	Avatar       int64     `json:"avatar" gorm:"column:avatar"`
	StudentId    string    `json:"student_id"`
	School       string    `json:"school"`
	Signature    string    `json:"signature"`
	StuId        string    `json:"stuId" gorm:"column:stu_id"`
	PublicKey    string    `gorm:"column:public_key;not null"`
	PrivateKey   string    `gorm:"column:private_key;not null"`
	FollowingNum int16     `json:"followingNum" gorm:"following_num"`
	FansNum      int16     `json:"fansNum" gorm:"fans_num"`
	CreateTime   time.Time `gorm:"column:create_time;autoCreateTime:milli"`
	UpdateTime   int64     `gorm:"column:update_time;autoUpdateTime:milli"`
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

func (user *User) FindByUsername() error {
	err := orm.DB.Where("username = ?", user.UserName).First(user).Error
	if err != nil && err == gorm.ErrRecordNotFound {
		return errors.New("user not exist")
	}
	return nil
}

func (user *User) FindById() error {
	err := orm.DB.Where("id = ?", user.Id).First(user).Error
	if err != nil && err == gorm.ErrRecordNotFound {
		return errors.New("user not exist")
	}
	return nil
}

func (user *User) Exists() bool {
	err := orm.DB.Model(user).First(user, "id=?", user.Id).Error
	return err != gorm.ErrRecordNotFound
}

func (user *User) ChangeUserAvatar(imageId int64) error {
	if !user.Exists() {
		return errors.New("user not exist")
	}
	return orm.DB.Model(user).Update("avatar", imageId).Error
}

func (user *User) ChangeUserProfile(attr string, value string) error {
	if !user.Exists() {
		return errors.New("user not exist")
	}
	if orm.DB.Model(user).Update(attr, value).RowsAffected == 0 {
		return errors.New("no such attribute")
	}
	return nil
}

func GetLikedUsers(articleId int64, pageSize int, pageNum int) (result []User, err error) {
	err = orm.DB.Raw("select * from user where id in (?) limit ?,?", orm.DB.Raw("select user_id from user_like_articles where article_id = ? ", articleId), pageSize*(pageNum-1), pageSize).Scan(&result).Error
	return
}

func GetFans(userId int64, pageSize int, pageNum int) (result []User, err error) {
	err = orm.DB.Raw("select * from user where id in (?) limit ?,?", orm.DB.Raw("select user_id from follows where following_id = ? ", userId), pageSize*(pageNum-1), pageSize).Scan(&result).Error
	return
}

func GetFollowing(userId int64, pageSize int, pageNum int) (result []User, err error) {
	err = orm.DB.Raw("select * from user where id in (?)  limit ?,?", orm.DB.Raw("select following_id from follows where user_id = ?", userId), pageSize*(pageNum-1), pageSize).Scan(&result).Error
	return
}

func SearchUser(key string, pageSize int, pageNum int) (result []User, err error) {
	//println("size:%d,num:%d",pageSize,pageNum)
	query := "%" + key + "%"
	err = orm.DB.Raw("select * from user where username like ? or nickname like? limit ?,?", query, query, pageSize*(pageNum-1), pageSize).Scan(&result).Error
	return
}

func CountUserNum() (result int64, err error) {
	err = orm.DB.Raw("select count(*) from user").Scan(&result).Error
	return
}

func (user *User) DailyUpdate(appVersion int64, stuId string) (err error) {
	return orm.DB.Model(user).Update("id", user.Id).Update("app_version", appVersion).Update("stu_id", stuId).Error
}
