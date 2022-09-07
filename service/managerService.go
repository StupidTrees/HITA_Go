package service

import (
	repo "hita/repository"
	"hita/utils/api"
	"strconv"
)

type CheckUpdateReq struct {
	VersionCode int64  `form:"versionCode" json:"versionCode"`
	Client      string `form:"client" json:"client"`
	Id          string `form:"id" json:"id"`
}
type CheckUpdateResult struct {
	ShouldUpdate      bool   `form:"shouldUpdate" json:"shouldUpdate"`
	LatestVersionCode int64  `form:"latestVersionCode" json:"latestVersionCode"`
	LatestVersionName string `form:"latestVersionName" json:"latestVersionName"`
	LatestUrl         string `form:"latestUrl" json:"latestUrl"`
	UpdateLog         string `form:"updateLog" json:"updateLog"`
}

func (req *CheckUpdateReq) CheckUpdate(userId int64) (data CheckUpdateResult, code int, er error) {
	data = CheckUpdateResult{}
	if req.Client == "android" {
		info := repo.Info{
			Key: "latest_version_code_android",
		}
		er = info.Get()
		user := repo.User{Id: userId}
		err := user.DailyUpdate(req.VersionCode, req.Id)
		if err != nil {
		}
		data.LatestVersionCode, er = strconv.ParseInt(info.Value, 10, 64)
		data.ShouldUpdate = data.LatestVersionCode > req.VersionCode
		info = repo.Info{
			Key: "latest_version_name",
		}
		er = info.Get()
		data.LatestVersionName = info.Value
		info = repo.Info{
			Key: "latest_version_url_android",
		}
		er = info.Get()
		data.LatestUrl = info.Value
		info = repo.Info{
			Key: "latest_version_log_android",
		}
		er = info.Get()
		data.UpdateLog = info.Value
	}
	if er != nil {
		code = api.CodeOtherError
	} else {
		code = api.CodeSuccess
	}
	return
}
