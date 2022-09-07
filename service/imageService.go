package service

import (
	"hita/config"
	"hita/repository"
	"hita/utils/logger"
	"io/ioutil"
	"path"
)

func GetImage(imageId int64) (data []byte, err error) {
	image := repository.Image{
		Id: imageId,
	}
	err = image.Find()
	if err != nil {
		return nil, err
	}
	var pt string
	switch image.Type {
	case "AVATAR":
		pt = config.AvatarPath
		break
	default:
		pt = config.ArticleImagePath
		break
	}
	fullPath := path.Join(logger.GetCurrentPath(), "..") + "/" + pt + image.Filename
	if image.Sensitive {
		fullPath = repository.GetSensitivePlaceholderPath()
	}
	data, err = ioutil.ReadFile(fullPath)
	return
}
