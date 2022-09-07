package repository

import (
	"hita/config"
	"hita/utils/logger"
	orm "hita/utils/mysql"
	"os"
	"path"
)

type Image struct {
	Id        int64  `json:"id" gorm:"PRIMARY_KEY;AUTO_INCREMENT"`
	Filename  string `json:"filename"`
	Type      string `json:"type" gorm:"type:enum('AVATAR','POST','OTHER');default:'OTHER';not null"`
	Sensitive bool   `gorm:"not null;default:0"`
}

func (Image) TableName() string {
	return "image"
}

func (i *Image) Create() error {
	result := orm.DB.Create(i)
	if result.Error != nil {
		logger.Errorln(result.Error)
		return result.Error
	}
	return nil
}
func GetAvatarPath(filename string) string {
	return path.Join(logger.GetCurrentPath(), "..") + "/" + config.AvatarPath + filename
}
func GetArticleImagePath(filename string) string {
	return path.Join(logger.GetCurrentPath(), "..") + "/" + config.ArticleImagePath + filename
}

func GetSensitivePlaceholderPath() string {
	return path.Join(logger.GetCurrentPath(), "..") + "/" + config.SensitivePlaceholderPath
}
func (i *Image) Delete() error {
	err := i.Find()
	if err == nil {
		var filePath string
		switch i.Type {
		case "AVATAR":
			{
				filePath = GetAvatarPath(i.Filename)
			}
		case "POST":
			{
				filePath = GetArticleImagePath(i.Filename)
			}
		}
		_ = os.Remove(filePath) //删除原先文件
		return orm.DB.Where("id=?", i.Id).Delete(i).Error
	} else {
		return err
	}
}
func (i *Image) Find() error {
	return orm.DB.Where("id=?", i.Id).Find(i).Error
}
func (i *Image) Update() error {
	return orm.DB.Model(i).Update("sensitive", i.Sensitive).Error
}
